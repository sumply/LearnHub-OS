#include "../include/responses.hxx"
#include <nlohmann/json_fwd.hpp>
#include <nlohmann/json.hpp>

namespace service_a {

response_builder::
response_builder(const http_request& request)
	: request(request) {}

http_response 
response_builder::
build_base_response(http::status status, const std::string& body) {
	http_response response(status, request.version());
	response.set(http::field::server, BOOST_BEAST_VERSION_STRING);
	response.set(http::field::content_type, "application/json");
	response.set(http::field::access_control_allow_origin, "*");
	response.set(http::field::access_control_allow_methods, "POST, GET, OPTIONS");
	response.set(http::field::access_control_allow_headers, "Content-Type, Authorization");
	response.set(http::field::access_control_max_age, "86400");
	response.keep_alive(request.keep_alive());
	response.body() = body;
	response.prepare_payload();

	return response;
}

http_response 
response_builder::
body_empty() {
	nlohmann::ordered_json text {
		{"error", "The request body is empty!"}
	};
	return build_base_response(http::status::bad_request, text.dump());
}

http_response 
response_builder::
invalid_json() {
	nlohmann::ordered_json text {
		{"error", "The request body is invalid!"},
		{"example", {
			{"login", "something"}, 
			{"password", "something"}}}
	};
	return build_base_response(http::status::bad_request, text.dump());
}

http_response 
response_builder::
method_not_allowed() {
	nlohmann::ordered_json text {
		{"error", "Use the GET method for additional information, or POST."}
	};
	return build_base_response(http::status::method_not_allowed, text.dump());
}

http_response 
response_builder::
missing_or_empty_key(const std::string& key) {
	nlohmann::ordered_json text {
		{"error", "Invalid JSON: the '" + key + "' is missing or empty"}
	};
	return build_base_response(http::status::bad_request, text.dump());
}

http_response 
response_builder::
unauthorized() {
	nlohmann::ordered_json text {
		{"error", "Authorization failed"}
	};
	return build_base_response(http::status::unauthorized, text.dump());
}

http_response
response_builder::
auth_jwt(const jwt_config::auth_tokens& tokens) {
	nlohmann::ordered_json text {
		{"access", tokens.access},
		{"refresh", tokens.refresh}
	};
	return build_base_response(http::status::ok, text.dump());
}

http_response 
response_builder::
access_jwt(const std::string& accessToken) {
	nlohmann::ordered_json text {{"access", accessToken}};
	return build_base_response(http::status::ok, text.dump());
}

http_response 
response_builder::
invalid_jwt_token() {
	nlohmann::ordered_json text {
		{"error", "The Authorization field is missing or the JWT token is invalid."}
	};
	return build_base_response(http::status::unauthorized, text.dump());
}

http_response 
response_builder::
unsupported_content_type() {
	nlohmann::ordered_json text {
		{"error", "Unsupported content type!"},
		{"Support"}, {
			{"application/json"}
		}
	};
	return build_base_response(http::status::bad_request, text.dump());
}

http_response
response_builder::
internal_server_error() {
	return build_base_response(http::status::internal_server_error,
		nlohmann::json{"error", "An error has occurred on the server side."});
}


http_response 
response_builder::
build_base_response(const http_request& request, http::status status, const std::string& body) {
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
invalid_jwt_token(const http_request& request) {
	nlohmann::ordered_json text {
		{"error", "The Authorization field is missing or the JWT token is invalid."}
	};
	return build_base_response(request, http::status::unauthorized, text.dump());
}


}
