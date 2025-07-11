#include "../include/protocol.hxx"

namespace protocol {

std::string to_string(target target) {
    switch (target) {
        case target::registration:
            return "/registration";
        case target::authorization:
            return "/authorization";
        case target::refresh_access_token:
            return "/refresh_access_token";
        case target::refresh_refresh_token:
            return "/refresh_refresh_token";
        case target::unknown:
            return "unknown";
    }
}

target to_target(const std::string& target) {
    if (target == "/registration")
        return target::registration;
    if (target == "/authorization")
        return target::authorization;
    if (target == "/refresh_access_token")
        return target::refresh_access_token;
    if (target == "/refresh_refresh_token")
        return target::refresh_refresh_token;
    return target::unknown;
}

}
