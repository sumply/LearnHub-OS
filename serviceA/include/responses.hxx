#pragma once

#include <boost/beast.hpp>
#include <boost/asio.hpp>
#include "../include/beast.hxx"

namespace service_a {

class response_builder {
public:
	explicit response_builder(http_request const& request);

	http_response buildBaseResponse(
			boost::beast::http::status status,
			const std::string& body);
	http_response bodyEmpty();
	http_response invalidJSON();
	http_response methodNotAllowed();
	http_response missingOrEmptyKey(const std::string& key);
	http_response unauthorized();
	//http_response authJWT(const jwtconfig::AuthTokens& tokens);
	http_response invalidJWTToken();
	http_response accessJWT(const std::string& accessToken);
	http_response unsupportedContentType();

	static http_response buildBaseResponse(
			const http_request& request,
			boost::beast::http::status status,
			const std::string& body);
	static http_response invalidJWTToken(const http_request& request);
private:
	const http_request& request;
};

}
