#pragma once

#include <boost/beast.hpp>
#include <boost/asio.hpp>
#include <boost/beast/core/error.hpp>
#include <memory>

namespace service_a {

class listener : public std::enable_shared_from_this<listener> {
public:
    listener(
        boost::asio::io_context &ioc,
        boost::asio::ip::tcp::endpoint endpoint
    );
    void run();
private:
    boost::asio::io_context& ioc_;
    boost::asio::ip::tcp::acceptor acceptor_;
    void on_accept(boost::beast::error_code ec, boost::asio::ip::tcp::socket socket);
    void fail(boost::beast::error_code ec, std::string_view what);
};

}
