#pragma once

#include <string>

namespace protocol {

enum class target {
    registration,
    authorization,
    get_user,
    refresh_refresh_token,
    refresh_access_token,
    unknown
};

std::string to_string(target target);
target to_target(const std::string& target);

}

