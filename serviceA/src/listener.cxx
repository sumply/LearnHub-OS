#include "../include/listener.hxx"
#include <boost/beast/core/error.hpp>
#include <iostream>
#include <format>
#include "../include/http_session.hxx"
#include "../include/logging.hxx"

namespace service_a {

listener::
listener(
    boost::asio::io_context& ioc,
    boost::asio::ip::tcp::endpoint endpoint)
    : ioc_(ioc)
    , acceptor_(ioc) {
        DEBUG_FUNC();
        boost::beast::error_code ec;

        ec = acceptor_.open(endpoint.protocol(), ec);
        if (ec) {
            fail(ec, "open");
            return;
        }
        ec = acceptor_.set_option(boost::asio::socket_base::reuse_address(true), ec);
        if (ec) {
            fail (ec, "set_option");
            return;
        }
        ec = acceptor_.bind(endpoint, ec);
        if (ec) {
            fail(ec, "bind");
            return;
        }
        ec = acceptor_.listen(
            boost::asio::socket_base::max_listen_connections, ec);
        if (ec) {
            fail(ec, "listen");
            return;
        }
    }

void
listener::
run() {
    DEBUG_FUNC();
    acceptor_.async_accept(
        boost::asio::make_strand(ioc_),
        boost::beast::bind_front_handler(
            &listener::on_accept,
            shared_from_this()));
}

void
listener::
fail(boost::beast::error_code ec, std::string_view what) {
    DEBUG_FUNC();
    if (ec == boost::asio::error::operation_aborted)
        return;
    std::cerr << std::format("FAIL: {}\n", what);
}

void
listener::
on_accept(boost::beast::error_code ec, boost::asio::ip::tcp::socket socket) {
    DEBUG_FUNC();
    if (ec)
        return fail(ec, "accept");
    run();
    std::make_shared<http_session>(std::move(socket))->run();
}

}
