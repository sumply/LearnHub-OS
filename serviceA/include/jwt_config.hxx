#pragma once

#include <jwt-cpp/jwt.h>
#include <jwt-cpp/traits/nlohmann-json/traits.h>
#include <boost/uuid.hpp>
#include <boost/uuid/random_generator.hpp>
#include <chrono>
#include "../include/beast.hxx"
#include "../include/logging.hxx"
#include "logging.hxx"
#include <optional>


namespace jwt_config {

using decoded_jwt = jwt::decoded_jwt<jwt::traits::nlohmann_json>;

class jwt_verifier {
public:
    jwt_verifier(const jwt_verifier&) = delete;
    jwt_verifier(jwt_verifier&&) = delete;
    jwt_verifier operator=(const jwt_verifier&) = delete;
    jwt_verifier operator=(jwt_verifier&&) = delete;
    static jwt_verifier& instance();
    bool verify(decoded_jwt const& decoded);
    bool verify(std::string_view token);
private:
    jwt_verifier();
    decltype(jwt::verify<jwt::traits::nlohmann_json>()) verifier_;
};

//using jwt_verifier = decltype(jwt::verify<jwt::traits::nlohmann_json>());
//using decoded_jwt = jwt::decoded_jwt<jwt::traits::nlohmann_json>;

struct auth_tokens {
    std::string access;
    std::string refresh;
};

std::string get_secret_key();
std::string make_access_token(const std::string& userID);
std::string make_refresh_token(const std::string& userID);
auth_tokens make_auth_tokens(const std::string& userID);
std::optional<std::string> get_bearer_token(const std::string& field_auth);
std::optional<std::string> get_bearer_token(const http_request& req);
std::optional<decoded_jwt> do_decode(std::string_view token);

template<typename Rep, typename Per>
std::string make_base_token(const std::string& userID, const std::chrono::duration<Rep, Per>& duration) {
    DEBUG_FUNC();
    auto now { std::chrono::system_clock::now() };
    auto exp { duration };

    boost::uuids::uuid jti { boost::uuids::random_generator()() };

    return jwt::create<jwt::traits::nlohmann_json>()
        .set_issuer("martin")
        .set_type("JWT")
        .set_subject(userID)
        .set_id(boost::uuids::to_string(jti))
        .set_issued_at(now)
        .set_expires_at(now + exp)
        .sign(jwt::algorithm::hs256{get_secret_key()});
}



}
