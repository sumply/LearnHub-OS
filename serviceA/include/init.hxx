#pragma once

#include "../include/beast.hxx"
#include <optional>

namespace service_a {

std::optional<ip::tcp::endpoint> endpoint_from_env(std::string_view env_address, std::string_view env_port);

}
