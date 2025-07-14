#include "../include/jwt_config.hxx"
#include "../include/dotenv.hxx"
#include <regex>

namespace jwt_config {

std::string get_secret_key() {
    DEBUG_FUNC();
    static const std::string secret_key = []{
        auto key = dotenv::getenv("JWT_SECRET_KEY");
        if (!key) {
            std::cerr << "JWT_SECRET_KEY not found in env.\n";
            std::exit(EXIT_FAILURE);
        }
        return key.value();
    }();
    return secret_key;
}

std::string make_access_token(const std::string& userID) {
    return make_base_token(userID, std::chrono::seconds(900));
}

std::string make_refresh_token(const std::string& userID) {
    return make_base_token(userID, std::chrono::seconds(2592000/*30 days*/));
}

auth_tokens make_auth_tokens(const std::string& userID) {
    return {
        .access { make_access_token(userID) },
        .refresh { make_refresh_token(userID) }
    };
}

jwt_verifier::
jwt_verifier()
    : verifier_(jwt::verify<jwt::traits::nlohmann_json>()
        .allow_algorithm(jwt::algorithm::hs256{get_secret_key()})
        .with_issuer("LearnHub-OS")
        .with_type("JWT")) {
            DEBUG_FUNC();
        }

jwt_verifier&
jwt_verifier::
instance() {
    DEBUG_FUNC();
    static jwt_verifier inst;
    return inst;
}

bool
jwt_verifier::
verify(decoded_jwt const& decoded) {
    DEBUG_FUNC();
    try {
        verifier_.verify(decoded);
        return true;
    } catch (const jwt::error::token_verification_exception& ec)  {
        return false;
    }
}

bool
jwt_verifier::
verify(std::string_view token) {
    DEBUG_FUNC();
    auto decoded = do_decode(token);
    if (!decoded)
        return false;
    return verify(decoded.value());
}

std::optional<std::string> get_bearer_token(const std::string& fieldAuthorization) {
    DEBUG_FUNC();
    std::regex regex {R"(Bearer\s+([A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+))"};
    std::smatch match;
    if (std::regex_search(fieldAuthorization, match, regex))
        return match[1].str();
    return std::nullopt;
}

std::optional<std::string> get_bearer_token(const http_request& req) {
    DEBUG_FUNC();
    auto token { req.find("Authorization") };
    if (token == req.end())
        throw std::nullopt;
    return get_bearer_token(token->value());
}

std::optional<decoded_jwt> do_decode(std::string_view token) {
    DEBUG_FUNC();
    try {
        return jwt::decode<jwt::traits::nlohmann_json>(std::string(token));
    } catch (const jwt::error::token_verification_exception& ec) {
        return std::nullopt;
    }
}

}
