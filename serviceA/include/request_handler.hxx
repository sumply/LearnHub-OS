#pragma once

#include "../include/beast.hxx"
#include "../include/responses.hxx"
#include <nlohmann/json_fwd.hpp>
#include <nlohmann/json.hpp>
#include "../include/jwt_config.hxx"
#include "jwt_config.hxx"

namespace service_a {

/**
 * @brief Handler http request for creates a response.
 * Supports GET, POST and WebSocket upgrade handling.
 * and return http response. It is maybe has been exception or correct answer.
 * @param request HTTP request from socket.
 * */
class request_handler {
public:
	request_handler(
		http_request&& request,
		jwt_config::jwt_verifier& verifier);
	~request_handler() = default;

	http_response handle_request();

	http_response handle_auth();
	http_response handle_get();
	http_response handle_post();
	http_response handle_options();
	http_response handle_reqistration();
private:
	http_request request_;
	response_builder response_builder_;
	jwt_config::jwt_verifier& jwt_verifier_;
};

}
