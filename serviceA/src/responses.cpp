#include "../include/responses.hxx"
#include <nlohmann/json_fwd.hpp>
#include <nlohmann/json.hpp>


namespace server::protocol::http {

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

response_builder::
response_builder(const http_request& request)
	: request(request) {}

http_response 
response_builder::
buildBaseResponse(http::status status, const std::string& body) {
	http_response response(status, request.version());
	response.set(http::field::server, BOOST_BEAST_VERSION_STRING);
	response.set(http::field::content_type, "application/json");
	response.keep_alive(request.keep_alive());
	response.body() = body;
	response.prepare_payload();
	return response;
}

http_response 
response_builder::
bodyEmpty() {
	nlohmann::ordered_json text {
		{"error", "The request body is empty!"}
	};
	return buildBaseResponse(http::status::bad_request, text.dump());
}

http_response 
response_builder::
invalidJSON() {
	nlohmann::ordered_json text {
		{"error", "The request body is invalid!"},
		{"example", {
			{"login", "something"}, 
			{"password", "something"}}}
	};
	return buildBaseResponse(http::status::bad_request, text.dump());
}

http_response 
response_builder::
methodNotAllowed() {
	nlohmann::ordered_json text {
		{"error", "Use the GET method for additional information, or POST."}
	};
	return buildBaseResponse(http::status::method_not_allowed, text.dump());
}

http_response 
response_builder::
missingOrEmptyKey(const std::string& key) {
	nlohmann::ordered_json text {
		{"error", "Invalid JSON: the '" + key + "' is missing or empty"}
	};
	return buildBaseResponse(http::status::bad_request, text.dump());
}

http_response 
response_builder::
unauthorized() {
	nlohmann::ordered_json text {
		{"error", "Authorization failed"}
	};
	return buildBaseResponse(http::status::unauthorized, text.dump());
}

/*http_response
response_builder::
authJWT(const jwtconfig::AuthTokens& tokens) {
	nlohmann::ordered_json text {
		{"access", tokens.access},
		{"refresh", tokens.refresh}
	};
	return buildBaseResponse(http::status::ok, text.dump());
}*/

http_response 
response_builder::
accessJWT(const std::string& accessToken) {
	nlohmann::ordered_json text {{"access", accessToken}};
	return buildBaseResponse(http::status::ok, text.dump());
}

http_response 
response_builder::
invalidJWTToken() {
	nlohmann::ordered_json text {
		{"error", "The Authorization field is missing or the JWT token is invalid."}
	};
	return buildBaseResponse(http::status::unauthorized, text.dump());
}

http_response 
response_builder::
unsupportedContentType() {
	nlohmann::ordered_json text {
		{"error", "Unsupported content type!"},
		{"Support"}, {
			{"application/json"}
		}
	};
	return buildBaseResponse(http::status::bad_request, text.dump());
}

http_response 
response_builder::
buildBaseResponse(const http_request& request, http::status status, const std::string& body) {
	http_response response(status, request.version());
	response.set(http::field::server, BOOST_BEAST_VERSION_STRING);
	response.set(http::field::content_type, "application/json");
	response.keep_alive(request.keep_alive());
	response.body() = body;
	response.prepare_payload();
	return response;
}

http_response 
response_builder::
invalidJWTToken(const http_request& request) {
	nlohmann::ordered_json text {
		{"error", "The Authorization field is missing or the JWT token is invalid."}
	};
	return buildBaseResponse(request, http::status::unauthorized, text.dump());
}

}
