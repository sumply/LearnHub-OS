#pragma once

#include <nlohmann/json.hpp>
#include <optional>

namespace json_config {

std::optional<nlohmann::json> parse(std::string_view text);
template<typename T = std::string>
std::optional<T> find(nlohmann::json const& json, std::string_view key) {
    auto it = json.find(key);
    if (it != json.end() && !it->is_null()) {
        try {
            return it->get<T>();
        } catch ( ... ) {
            return std::nullopt;
        }
    }
    return std::nullopt;
}

}
