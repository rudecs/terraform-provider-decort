/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration 
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates. 
*/

package decort

import (

	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	// "time"

	"github.com/dgrijalva/jwt-go"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/terraform"

)

// enumerated constants that define authentication modes 
const (
	MODE_UNDEF   = iota // this is the invalid mode - it should never be seen
	MODE_LEGACY  = iota
	MODE_OAUTH2  = iota
	MODE_JWT     = iota
)

type ControllerCfg struct {
	controller_url   string  // always required
	auth_mode_code   int     // always required
	auth_mode_txt    string  // always required, it is a text representation of auth mode
	legacy_user      string  // required for legacy mode
	legacy_password  string  // required for legacy mode
	legacy_sid       string  // obtained from DECORT controller on successful login in legacy mode
	jwt              string  // obtained from Outh2 provider on successful login in oauth2 mode, required in jwt mode
	app_id           string  // required for oauth2 mode
	app_secret       string  // required for oauth2 mode
	oauth2_url       string  // always required
	decort_username    string  // assigned to either legacy_user (legacy mode) or Oauth2 user (oauth2 mode) upon successful verification
	cc_client        *http.Client // assigned when all initial check successfully passed
}

func ControllerConfigure(d *schema.ResourceData) (*ControllerCfg, error) {
	// This function first will check that all required provider parameters for the 
	// selected authenticator mode are set correctly and initialize ControllerCfg structure
	// based on the provided parameters.
	//
	// Next, it will check for validity of supplied credentials by initiating connection to the specified 
	// DECORT controller URL and, if succeeded, completes ControllerCfg structure with the rest of computed 
	// parameters (e.g. JWT, session ID and Oauth2 user name).
	//
	// The structure created by this function should be used with subsequent calls to decortAPICall() method, 
	// which is a DECORT authentication mode aware wrapper around standard HTTP requests.

	ret_config := &ControllerCfg{
		controller_url:  d.Get("controller_url").(string),
		auth_mode_code:  MODE_UNDEF,
		legacy_user:     d.Get("user").(string),
		legacy_password: d.Get("password").(string),
		legacy_sid:      "",
		jwt:             d.Get("jwt").(string),
		app_id:          d.Get("app_id").(string),
		app_secret:      d.Get("app_secret").(string),
		oauth2_url:      d.Get("oauth2_url").(string),
		decort_username:   "",
	}

	var allow_unverified_ssl bool
	allow_unverified_ssl = d.Get("allow_unverified_ssl").(bool)

	if ret_config.controller_url == "" {
		return nil, fmt.Errorf("Empty DECORT cloud controller URL provided.")
	}

	// this should have already been done by StateFunc defined in Schema, but we want to be sure
	ret_config.auth_mode_txt = strings.ToLower(d.Get("authenticator").(string))

	switch ret_config.auth_mode_txt {
	case "jwt":
		if ret_config.jwt == "" {
			return nil, fmt.Errorf("Authenticator mode 'jwt' specified but no JWT provided.")
		}
		ret_config.auth_mode_code = MODE_JWT
	case "oauth2":
		if ret_config.oauth2_url == "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no OAuth2 URL provided.")
		}
		if ret_config.app_id == "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no Application ID provided.")
		}
		if ret_config.app_secret == "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no Secret ID provided.")
		}
		ret_config.auth_mode_code = MODE_OAUTH2
	case "legacy":
		//
		ret_config.legacy_user = d.Get("user").(string)
		if ret_config.legacy_user == "" {
			return nil, fmt.Errorf("Authenticator mode 'legacy' specified but no user provided.")
		}
		ret_config.legacy_password = d.Get("password").(string)
		if ret_config.legacy_password == "" {
			return nil, fmt.Errorf("Authenticator mode 'legacy' specified but no password provided.")
		}
		ret_config.auth_mode_code = MODE_LEGACY
	default:
		return nil, fmt.Errorf("Unknown authenticator mode %q provided.", ret_config.auth_mode_txt)
	}

	if allow_unverified_ssl {
		log.Printf("ControllerConfigure: allow_unverified_ssl is set - will not check certificates!")
		transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
		ret_config.cc_client = &http.Client{
			Transport: transCfg,
			Timeout: Timeout180s,
		}
	} else {
		ret_config.cc_client = &http.Client{
			Timeout: Timeout180s, // time.Second * 30,
		}
	}

	switch ret_config.auth_mode_code {
	case MODE_LEGACY:
		ok, err := ret_config.validateLegacyUser() 
		if !ok {
			return nil, err
		}
		ret_config.decort_username = ret_config.legacy_user
	case MODE_JWT:
		//
		ok, err := ret_config.validateJWT("")
		if !ok {
			return nil, err
		}
	case MODE_OAUTH2:
		// on success getOAuth2JWT will set config.jwt to the obtained JWT, so there is no 
		// need to set it once again here
		_, err := ret_config.getOAuth2JWT()
		if err != nil {
			return nil, err
		}
		// we are not verifying the JWT when parsing because actual verification is done on the 
		// OVC controller side. Here we do parsing solely to extract Oauth2 user name (claim "user")
		// and JWT issuer name (claim "iss")
		parser := jwt.Parser{}
		token, _, err := parser.ParseUnverified(ret_config.jwt, jwt.MapClaims{})
		if err != nil {
			return nil, err
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			var tbuf bytes.Buffer
			tbuf.WriteString(claims["username"].(string))
			tbuf.WriteString("@")
			tbuf.WriteString(claims["iss"].(string))
			ret_config.decort_username =  tbuf.String() 
		} else {
			return nil, fmt.Errorf("Failed to extract user and iss fields from JWT token in oauth2 mode.")
		}
	default:
		// FYI, this should never happen due to all above checks, but we want to be fool proof
		return nil, fmt.Errorf("Unknown authenticator mode code %d provided.", ret_config.auth_mode_code)
	}

	// All checks passed successfully, credentials corresponding to the selected authenticator mode
	// obtained and validated.
	return ret_config, nil
}

func (config *ControllerCfg) getDecsUsername() (string) {
	return config.decort_username
}

func (config *ControllerCfg) getOAuth2JWT() (string, error) {
	// 	Obtain JWT from the Oauth2 provider using application ID and application secret provided in config.
	if config.auth_mode_code == MODE_UNDEF {
		return "", fmt.Errorf("getOAuth2JWT method called for undefined authorization mode.")
	}
	if config.auth_mode_code != MODE_OAUTH2 {
		return "", fmt.Errorf("getOAuth2JWT method called for incompatible authorization mode %q.", config.auth_mode_txt)
	}

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", config.app_id)
	params.Add("client_secret", config.app_secret)
	params.Add("response_type", "id_token")
	params.Add("validity", "3600")
	params_str := params.Encode()

	req, err := http.NewRequest("POST", config.oauth2_url + "/v1/oauth/access_token", strings.NewReader(params_str))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(params_str)))

	resp, err := config.cc_client.Do(req)	
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		// fmt.Println("response Status:", resp.Status)
		// fmt.Println("response Headers:", resp.Header)
		// fmt.Println("response Headers:", req.URL)
		return "", fmt.Errorf("getOauth2JWT: unexpected status code %d when obtaining JWT from %q for APP_ID %q, request Body %q", 
		resp.StatusCode, req.URL, config.app_id, params_str)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
 
	// validation successful - store JWT in the corresponding field of the ControllerCfg structure
	config.jwt = strings.TrimSpace(string(responseData))

	return config.jwt, nil
}

func (config *ControllerCfg) validateJWT(jwt string) (bool, error) {
	/*
	Validate JWT against DECORT controller. JWT can be supplied as argument to this method. If empty string supplied as
	argument, JWT will be taken from config attribute. 
	DECORT controller URL will always be taken from the config attribute assigned at instantiation.
    Validation is accomplished by attempting API call that lists accounts for the invoking user.
	*/
	if jwt == "" {
		if config.jwt == "" {
			return false, fmt.Errorf("validateJWT method called, but no meaningful JWT provided.")
		}
		jwt = config.jwt
	}

	if config.oauth2_url == "" {
		return false, fmt.Errorf("validateJWT method called, but no OAuth2 URL provided.")
	}

	req, err := http.NewRequest("POST", config.controller_url + "/restmachine/cloudapi/accounts/list", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Length", strconv.Itoa(0))
	
	resp, err := config.cc_client.Do(req)	
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("validateJWT: unexpected status code %d when validating JWT against %q.", 
		resp.StatusCode, req.URL)
	}
	defer resp.Body.Close()

	return true, nil
}

func (config *ControllerCfg) validateLegacyUser() (bool, error) {
	/*
	Validate legacy user by obtaining a session key, which will be used for authenticating subsequent API calls
    to DECORT controller.
    If successful, the session key is stored in config.legacy_sid and true is returned. If unsuccessful for any
    reason, the method will return false and error.
	*/
	if config.auth_mode_code == MODE_UNDEF {
		return false, fmt.Errorf("validateLegacyUser method called for undefined authorization mode.")
	}
	if config.auth_mode_code != MODE_LEGACY {
		return false, fmt.Errorf("validateLegacyUser method called for incompatible authorization mode %q.", config.auth_mode_txt)
	}

	params := url.Values{}
	params.Add("username", config.legacy_user)
	params.Add("password", config.legacy_password)
	params_str := params.Encode()

	req, err := http.NewRequest("POST", config.controller_url + "/restmachine/cloudapi/users/authenticate", strings.NewReader(params_str))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(params_str)))

	resp, err := config.cc_client.Do(req)	
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("validateLegacyUser: unexpected status code %d when validating legacy user %q against %q.", 
		resp.StatusCode, config.legacy_user, config.controller_url)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
 
	// validation successful - keep session ID for future use
	config.legacy_sid = string(responseData)

	return true, nil
}

func (config *ControllerCfg) decortAPICall(method string, api_name string, url_values *url.Values) (json_resp string, err error) {
	// This is a convenience wrapper around standard HTTP request methods that is aware of the 
	// authorization mode for which the provider was initialized and compiles request accordingly.

	if config.cc_client == nil {
		// this should never happen if ClientConfig was properly called prior to decortAPICall 
		return "", fmt.Errorf("decortAPICall method called with unconfigured DECORT cloud controller HTTP client.")
	}

	// Example: to create api_params, one would generally do the following:
	// 
	// data := []byte(`{"machineId": "2638"}`)
	// api_params := bytes.NewBuffer(data))
	//
	// Or:
	//
	// params := url.Values{}
	// params.Add("machineId", "2638")
	// params.Add("username", "u")
	// params.Add("password", "b")
	// req, _ := http.NewRequest(method, url, strings.NewReader(params.Encode()))
	//

	if config.auth_mode_code == MODE_UNDEF {
		return "", fmt.Errorf("decortAPICall method called for unknown authorization mode.")
	}

	if config.auth_mode_code == MODE_LEGACY {
		url_values.Add("authkey", config.legacy_sid)
	}
	params_str := url_values.Encode()

	req, err := http.NewRequest(method, config.controller_url + api_name, strings.NewReader(params_str))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(params_str)))

	if config.auth_mode_code == MODE_OAUTH2 || config.auth_mode_code == MODE_JWT {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", config.jwt))
	} 
	
	resp, err := config.cc_client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
		tmp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		} 
		json_resp := Jo2JSON(string(tmp_body))
		log.Printf("decortAPICall:\n %s", json_resp)
		return json_resp, nil
	} else {
		return "", fmt.Errorf("decortAPICall: unexpected status code %d when calling API %q with request Body %q", 
		resp.StatusCode, req.URL, params_str)
	}
	
	/*
	if resp.StatusCode == StatusServiceUnavailable {
        return nil, fmt.Errorf("decortAPICall method called for incompatible authorization mode %q.", config.auth_mode_txt)
	}
	*/

	return "", err
}

