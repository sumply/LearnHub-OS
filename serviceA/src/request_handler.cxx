#include "../include/request_handler.hxx"
#include "../include/protocol.hxx"
#include "../include/json_config.hxx"
#include "../include/logging.hxx"
#include <iostream>

namespace service_a {

request_handler::
request_handler(
    http_request&& request,
    jwt_config::jwt_verifier& verifier)
    : request_(std::move(request))
    , jwt_verifier_(verifier)
    , response_builder_(request) {
        DEBUG_FUNC();
}

http_response
request_handler::
handle_request() {
    DEBUG_FUNC();
    switch (request_.method()) {
        case http::verb::get:
            return handle_get();
        case http::verb::post:
            if (request_[http::field::content_type] != "application/json")
                return response_builder_.unsupported_content_type();
            return handle_post();
        default:
            return response_builder_.method_not_allowed();
    }
}

http_response
request_handler::
handle_auth() {
    DEBUG_FUNC();
    auto json = json_config::parse(request_.body());
    if (!json)
        return response_builder_.invalid_json();

    auto email = json_config::find(json, "email");
    if (!email)
        return response_builder_.missing_or_empty_key("email");

    auto password = json_config::find(json, "password");
    if (!password)
        return response_builder_.missing_or_empty_key("password");

   // FIXME: Request database userID.
    int userID = 10;

    auto tokens = jwt_config::make_auth_tokens(std::to_string(userID));
    return response_builder_.auth_jwt(tokens);
}

http_response
request_handler::
handle_reqistration() {
    DEBUG_FUNC();
    auto json = json_config::parse(request_.body());
    if (!json)
        return response_builder_.invalid_json();

    auto email = json_config::find(json, "email");
    if (!email)
        return response_builder_.missing_or_empty_key("email");

    auto password = json_config::find(json, "password");
    if (!password)
        return response_builder_.missing_or_empty_key("password");

    auto name = json_config::find(json, "name");
    if (!password)
        return response_builder_.missing_or_empty_key("name");

    auto surname = json_config::find(json, "surname");
    if (!surname)
        return response_builder_.missing_or_empty_key("surname");

    auto group = json_config::find(json, "group");
    if (!group)
        return response_builder_.missing_or_empty_key("group");

    auto role = json_config::find(json, "role");
    if (!role)
        return response_builder_.missing_or_empty_key("role");

    //FIXME: Request database.

    return response_builder_.build_base_response(
        request_, boost::beast::http::status::ok,
        "Авторизация прошла успешно!");
}


http_response
request_handler::
handle_post() {
    DEBUG_FUNC();
    if (request_.body().empty())
        return response_builder_.body_empty();
    switch (protocol::to_target(request_.target())) {
        case protocol::target::authorization:
            return handle_auth();
        case protocol::target::registration:
            return handle_reqistration();
    }
}

http_response
request_handler::
handle_get() {
    DEBUG_FUNC();
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

    return response_builder_.build_base_response(http::status::ok, text.dump());
}

}
