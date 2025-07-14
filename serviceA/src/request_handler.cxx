#include "../include/request_handler.hxx"
#include "../include/protocol.hxx"
#include "../include/json_config.hxx"
#include "../include/logging.hxx"
#include "../include/http_client.hxx"
#include "../include/dotenv.hxx"

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
        case http::verb::options:
            return handle_options();
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
handle_options() {
    return response_builder_.build_base_response(http::status::no_content, "");
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
    http_client client(http_client::method::POST, dotenv::getenv("SERVICEB_AUTH").value());
    auto text = nlohmann::json{
        {"email", email.value()},
        {"password", password.value()}
    }.dump();
    client.field("Content-type: application/json");
    client.body(text);
    client.show();
    long status = client.send();
    if (status == -1)
        return response_builder_.internal_server_error();
    if (status != 200)
        return response_builder_.unauthorized();

    json = json_config::parse(client.response());
    if (!json)
        return response_builder_.invalid_json();

    auto user_id = json_config::find(json, "id");
    if (!user_id)
        return response_builder_.missing_or_empty_key("id");
    auto tokens = jwt_config::make_auth_tokens(user_id.value());
    return response_builder_.auth_jwt(tokens);
}

http_response
request_handler::
handle_reqistration() {
    DEBUG_FUNC();
    auto json = json_config::parse(request_.body());
    if (!json)
        return response_builder_.invalid_json();

    auto firstname = json_config::find(json, "firstname");
    if (!firstname)
        return response_builder_.missing_or_empty_key("firstname");
    auto secondname = json_config::find(json, "secondname");
    if (!secondname)
        return response_builder_.missing_or_empty_key("secondname");
    auto lastname = json_config::find(json, "lastname");
    if (!lastname)
        return response_builder_.missing_or_empty_key("lastname");
    auto email = json_config::find(json, "email");
    if (!email)
        return response_builder_.missing_or_empty_key("email");
    auto password = json_config::find(json, "password");
    if (!password)
        return response_builder_.missing_or_empty_key("password");
    auto role = json_config::find(json, "role");
    if (!role)
        return response_builder_.missing_or_empty_key("role");
    auto group = json_config::find(json, "group");
    if (!group)
        return response_builder_.missing_or_empty_key("group");

    auto service_response = post_registration({
        .firstname = firstname.value(),
        .secondname = secondname.value(),
        .lastname = lastname.value(),
        .email = email.value(),
        .password = password.value(),
        .group = group.value(),
        .role = role.value()
    });

    return response_builder_.build_base_response(
        request_, boost::beast::http::status::ok,
        service_response.dump());
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
        case protocol::target::get_user:
            return handle_get_user();
        case protocol::target::refresh_access_token:
        case protocol::target::refresh_refresh_token:
        case protocol::target::unknown:
            return response_builder_.build_base_response(http::status::bad_request, "{\"error\": \"Target\"}");
    }
}

http_response
request_handler::
handle_get_user() {
    DEBUG_FUNC();
    auto token = jwt_config::get_bearer_token(request_);
    if (!token)
        return response_builder_.invalid_jwt_token();

    auto decoded_token = jwt_config::do_decode(token.value());
    if (!decoded_token)
        return response_builder_.invalid_jwt_token();

    bool is_verified = jwt_verifier_.verify(decoded_token.value());
    if (!is_verified)
        return response_builder_.invalid_jwt_token();

    if (!decoded_token->has_payload_claim("sub"))
        return response_builder_.invalid_jwt_token();

    http_client client(http_client::method::POST, dotenv::getenv("SERVICEB_GET_USER").value());
    client.field("Content-type: application/json");
    auto text = nlohmann::json{{"id", decoded_token->get_payload_claim("sub").as_string()}}.dump();
    client.body(text);
    client.show();
    long status = client.send();
    if (status == -1)
        return response_builder_.internal_server_error();
    if (status != 200)
        return response_builder_.unauthorized();

    return response_builder_.build_base_response(http::status::ok, client.response());
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

nlohmann::json
request_handler::
post_registration(models::registration model) {
    DEBUG_FUNC();
    http_client client(http_client::method::POST, dotenv::getenv("SERVICEB_REG").value());
    client.field("Content-type: application/json");
    auto text = nlohmann::json{
        {"firstname", model.firstname},
        {"secondname", model.secondname},
        {"lastname", model.lastname},
        {"email", model.email},
        {"password", model.password},
        {"group", model.group},
        {"role", model.role}
    }.dump();
    client.body(text);
    client.show();
    long status = client.send();
    DEBUG_LOG(client.response());
    if (status == -1 || status != 200) {
        DEBUG_ERROR(client.error_message());
        return{};
    }
    auto json = json_config::parse(client.response());
    return json.value_or({});
}

nlohmann::json
request_handler::
post_authorizatoin(std::string_view email, std::string_view password) {
    DEBUG_FUNC();
    http_client client(http_client::method::POST, dotenv::getenv("SERVICEB_AUTH").value());
    client.field("Content-type: application/json");
    client.body(nlohmann::json{
        {"email", email},
        {"password", password}
    }.dump());
    client.show();
    long status = client.send();
    DEBUG_LOG(client.response());
    if (status == -1 || status != 200) {
        DEBUG_ERROR(client.error_message());
        return{};
    }
    auto json = json_config::parse(client.response());
    return json.value_or("");
}
}
