#include "../include/init.hxx"
#include "../include/dotenv.hxx"
#include <boost/asio/ip/address_v4.hpp>

namespace service_a {

std::optional<ip::tcp::endpoint> endpoint_from_env(std::string_view env_address, std::string_view env_port) {
    auto address_opt = dotenv::getenv(env_address);
    auto port_opt = dotenv::getenv(env_port);
    if (!address_opt || !port_opt)
        return std::nullopt;
    try {
        auto const address = asio::ip::make_address_v4(address_opt.value());
        auto const port = static_cast<unsigned short>(std::stoi(port_opt.value()));
        return ip::tcp::endpoint(address, port);
    } catch ( std::exception const& ec) {
        return std::nullopt;
    }
}

}
