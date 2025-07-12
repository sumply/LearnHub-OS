#pragma once

#include <boost/beast.hpp>
#include <boost/asio.hpp>
#include "../include/beast.hxx"
#include "../include/jwt_config.hxx"

namespace service_a {

class response_builder {
public:
	explicit response_builder(http_request const& request);

	http_response build_base_response(
			boost::beast::http::status status,
			const std::string& body);
	http_response body_empty();
	http_response invalid_json();
	http_response method_not_allowed();
	http_response missing_or_empty_key(const std::string& key);
	http_response unauthorized();
	http_response auth_jwt(const jwt_config::auth_tokens& tokens);
	http_response invalid_jwt_token();
	http_response access_jwt(const std::string& accessToken);
	http_response unsupported_content_type();

	static http_response build_base_response(
			const http_request& request,
			boost::beast::http::status status,
			const std::string& body);
	static http_response invalid_jwt_token(const http_request& request);
private:
	const http_request& request;
};

}
