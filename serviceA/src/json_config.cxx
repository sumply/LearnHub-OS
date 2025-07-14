#include "../include/json_config.hxx"
#include "../include/logging.hxx"

namespace json_config {

std::optional<nlohmann::json> parse(std::string_view text) {
    DEBUG_FUNC();
    auto parsed = nlohmann::json::parse(text, nullptr, false);
    if (parsed.is_discarded())
        return std::nullopt;
    return parsed;
}

}
