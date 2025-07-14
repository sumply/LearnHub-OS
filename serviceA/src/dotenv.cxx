#include "../include/dotenv.hxx"
#include "../include/logging.hxx"
#include <fstream>
#include <cstring>

namespace dotenv {

bool init(const std::filesystem::path& path) {
    DEBUG_FUNC();
    std::ifstream file(path);
    if (!file.is_open())
        return false;
    std::string line;
    while (std::getline(file, line)) {
        auto hash { line.find("#") };
        if (line.empty() || hash == 0)
            continue;
        auto equal { line.find("=") };
        if (equal == std::string::npos)
            return false;
        auto key { line.substr(0, equal) };
        auto value { line.substr(equal+1, std::min(hash, line.size()) - equal - 1) };
        setenv(key.c_str(), value.c_str(), 1);
    }
    file.close();
    return true;
}

std::optional<std::string> getenv(std::string_view key) {
    DEBUG_FUNC();
    const auto value {std::getenv(key.data())};
    if (value == nullptr || std::strlen(value) == 0)
        return std::nullopt;
    return value;
}

}
