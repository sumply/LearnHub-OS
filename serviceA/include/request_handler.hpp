#pragma once

#include "../include/beast.hxx"
#include "../include/responses.hxx"
#include <nlohmann/json_fwd.hpp>
#include <nlohmann/json.hpp>

namespace protocol::http {

	enum class Target {
		Registration,
		Authorization,
		RefreshRefreshToken,
		RefreshAccessToken,
		Unknown
	};

	std::string to_string(Target target);
	Target to_target(const std::string& target);

	std::string to_string(Target target) {
		switch (target) {
			case Target::Registration:
				return "/registration";
			case Target::Authorization:
				return "/authorization";
			case Target::RefreshAccessToken:
				return "/refresh_access_token";
			case Target::RefreshRefreshToken:
				return "/refresh_refresh_token";
			case Target::Unknown:
				return "unknown";
		}
	}

	Target to_target(const std::string& target) {
		if (target == "/registration")
			return Target::Registration;
		if (target == "/authorization")
			return Target::Authorization;
		if (target == "/refresh_access_token")
			return Target::RefreshAccessToken;
		if (target == "/refresh_refresh_token")
			return Target::RefreshRefreshToken;
		return Target::Unknown;
	}
}

namespace service_a {

/**
 * @brief Handler http request for creates a response.
 * Supports GET, POST and WebSocket upgrade handling.
 * and return http response. It is maybe has been exception or correct answer.
 * @param request HTTP request from socket.
 * */
class RequestHandler {
public:
	RequestHandler(
		http_request&& request);
	~RequestHandler() = default;

	http_response handleRequest();

	http_response handleAuthorization();
	//http_response handleRegistration();
	//http_response handleRefreshToken();
	http_response handlePost();
	http_response handleGet();
private:
	http_request request;
	//std::shared_ptr<database::PoolConnection<DBContext>> pool;
	//const jwtconfig::JWTVerifier& jwtVerifier;
	response_builder responseBuilder;
};


RequestHandler::
RequestHandler(
	http_request&& request)
	: request{std::move(request)}
	, responseBuilder{request} {
}

http_response
RequestHandler::
handleRequest() {
	switch (request.method()) {
	case http::verb::post:
		if (request[http::field::content_type] != "application/json")
			return responseBuilder.unsupportedContentType();
		return handlePost();
	case http::verb::get:
		return handleGet();
	default:
		return responseBuilder.methodNotAllowed();
	}
}

http_response
RequestHandler::
handleAuthorization() {
	try {
		nlohmann::json parsedBody = nlohmann::json::parse(request.body());

		//auto login { json::find<std::string>(parsedBody, "login") };
		//auto password { json::find<std::string>(parsedBody, "password") };

		//database::ScopedConnection con { *pool };

		//auto authResult { con->user().authorization({
		//		login,
		//		password}) };

		//if (!authResult)
		//	return responseBuilder.unauthorized();

		return responseBuilder.buildBaseResponse(http::status::bad_request, "Вот пароль и логин");
		//return responseBuilder.authJWT(jwtconfig::makeAuthTokens(std::to_string(authResult->id)));
	} catch (const nlohmann::json::parse_error& ec) {
		return responseBuilder.invalidJSON();
	//} catch (const JSONInvalidKey& ec) {
	//	return responseBuilder.invalidJSON();
	//}
}
}
/*
template<typename DBContext>
http_response 
RequestHandler<DBContext>::
handleRegistration() {
	try {
		nlohmann::json parsedBody = nlohmann::json::parse(request.body());

		auto nickname { json::find<std::string>(parsedBody, "nickname") };
		auto email { json::find<std::string>(parsedBody, "email") };
		auto password { json::find<std::string>(parsedBody, "password") };
database::ScopedConnection con { *pool };

		con->user().create({nickname, email, password});

		nlohmann::ordered_json text {
			{"what", "User was created!"}
		};
		return responseBuilder.buildBaseResponse(http::status::ok, text.dump());
	} catch (const nlohmann::json_abi_v3_12_0::detail::parse_error& ec) {
		return responseBuilder.invalidJSON();
	} catch (const JSONInvalidKey& ec) {
		return responseBuilder.invalidJSON();
	}
}

template<typename DBContext>
http_response 
RequestHandler<DBContext>::
handleRefreshToken() {
	try {
		auto decodedJWT { jwtconfig::doDecode(jwtconfig::getBearerToken(request)) };
		jwtconfig::doVerify(
			jwtVerifier, 
			decodedJWT);

		return responseBuilder.accessJWT(
			jwtconfig::makeAccessToken(decodedJWT.get_subject()));	
	} catch (const JSONInvalidKey& ec) {
		return responseBuilder.invalidJSON();
	} catch (const jwt::error::token_verification_exception& ec) {
		return responseBuilder.invalidJWTToken();
	} catch (const JWTAuthorizationFieldException& ec) {
		return responseBuilder.invalidJWTToken();
	}
}
*/
http_response
RequestHandler::
handlePost() {
	if (request.body().empty())
		return responseBuilder.bodyEmpty();

	using protocol::http::Target;
	switch (protocol::http::to_target(request.target())) {
//	case Target::Registration:
//		return handleRegistration();
	case Target::Authorization:
		return handleAuthorization();
//	case Target::RefreshAccessToken:
//		return handleRefreshToken();
//	case Target::RefreshRefreshToken:
//		break;
	case Target::Unknown:
		break;
	}
}

http_response
RequestHandler::
handleGet() {
	nlohmann::ordered_json text {
		{"status", 200},
		{"authorization_type", "JWT"},
		{"authorization_required", true},
		{"refresh_token_supported", false},
		{"allowed_methods", {"GET", "POST"}},
		{"allowed_protocols", {"wss", "https", "http", "ws"}},
		{"allowed_target", {
			protocol::http::to_string(protocol::http::Target::Registration),
			protocol::http::to_string(protocol::http::Target::Authorization),
			protocol::http::to_string(protocol::http::Target::RefreshAccessToken),
			protocol::http::to_string(protocol::http::Target::RefreshRefreshToken)}}
	};

	return responseBuilder.buildBaseResponse(http::status::ok, text.dump());
}

}
