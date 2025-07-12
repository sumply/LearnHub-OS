#pragma once

#include <string>
#include <filesystem>
#include <string_view>
#include <optional>

namespace dotenv {

/** Simple parsing of the .env file */
bool init(const std::filesystem::path& path);
/** Check that required environment keys are present */
std::optional<std::string> getenv(std::string_view key);

}
