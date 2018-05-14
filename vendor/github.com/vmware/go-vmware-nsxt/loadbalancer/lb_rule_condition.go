/*
 * NSX API
 *
 * VMware NSX REST API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package loadbalancer

// Match conditions are used to match application traffic passing through load balancers. Multiple match conditions can be specified in one load balancer rule, each match condition defines a criterion for application traffic. If inverse field is set to true, the match result of the condition is inverted. If more than one match condition is specified, match strategy determines if all conditions should match or any one condition should match for the load balancer rule to be considered a match. Currently only HTTP messages are supported by load balancer rules. Each load balancer rule is used at a specific phase of load balancer processing. Currently three phases are supported, HTTP_REQUEST_REWRITE, HTTP_FORWARDING and HTTP_RESPONSE_REWRITE. Each phase supports certain types of match conditions, supported match conditions in HTTP_REQUEST_REWRITE phase are: LbHttpRequestMethodCondition LbHttpRequestUriCondition LbHttpRequestUriArgumentsCondition LbHttpRequestVersionCondition LbHttpRequestHeaderCondition LbHttpRequestCookieCondition LbHttpRequestBodyCondition LbTcpHeaderCondition LbIpHeaderCondition Supported match conditions in HTTP_FORWARDING phase are: LbHttpRequestMethodCondition LbHttpRequestUriCondition LbHttpRequestVersionCondition LbHttpRequestHeaderCondition LbHttpRequestCookieCondition LbHttpRequestBodyCondition LbTcpHeaderCondition LbIpHeaderCondition Supported match condition in HTTP_RESPONSE_REWRITE phase is: LbHttpResponseHeaderCondition
type LbRuleCondition struct {

	// A flag to indicate whether reverse the match result of this condition
	Inverse bool `json:"inverse,omitempty"`

	// Type of load balancer rule condition
	Type_ string `json:"type"`
}
