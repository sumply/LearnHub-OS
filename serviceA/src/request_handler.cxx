#include "../include/request_handler.hxx"
#include "../include/protocol.hxx"

namespace service_a {

request_handler::
request_handler(
    http_request&& request)
    : request(std::move(request))
    , responseBuilder(request) {
}

http_response
request_handler::
handle_request() {
    switch (request.method()) {
        case http::verb::get:
            return handle_get();
        case http::verb::post:
            if (request[http::field::content_type] != "application/json")
                return responseBuilder.unsupported_content_type();
        default:
            return responseBuilder.method_not_allowed();
    }
}

http_response
request_handler::
handle_auth() {
    try {
        nlohmann::json parsedBody = nlohmann::json::parse(request.body());
        return responseBuilder.build_base_response(http::status::bad_request, "Вот пароль и логин");
    } catch (const nlohmann::json::parse_error& ec) {
        return responseBuilder.invalid_json();
    }
}

http_response
request_handler::
handle_reqistration() {

}

http_response
request_handler::
handle_post() {
    if (request.body().empty())
        return responseBuilder.body_empty();
    switch (protocol::to_target(request.target())) {
        case protocol::target::authorization:
            return handle_auth();
        case protocol::target::registration:
            return handle_reqistration();
    }
}

http_response
request_handler::
handle_get() {
    nlohmann::ordered_json text {
        {"status", 200},
        {"authorization_type", "JWT"},
        {"authorization_required", true},
        {"refresh_token_supported", false},
        {"allowed_methods", {"GET", "POST"}},
        {"allowed_protocols", {"http", "ws"}},
        {"allowed_target", {
            to_string(protocol::target::registration),
            to_string(protocol::target::authorization),
            to_string(protocol::target::refresh_access_token),
            to_string(protocol::target::refresh_refresh_token)}}
    };

    return responseBuilder.build_base_response(http::status::ok, text.dump());
}

}
