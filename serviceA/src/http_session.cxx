#include "../include/http_session.hxx"
#include <boost/beast/core/bind_handler.hpp>
#include <boost/beast/core/error.hpp>
#include <boost/beast/http/string_body_fwd.hpp>
#include "../include/request_handler.hxx"
#include <iostream>

namespace service_a {

http_session::
http_session(
        boost::asio::ip::tcp::socket socket)
    : stream_(std::move(socket))
{}

void
http_session::
run() {
    do_read();
}

void
http_session::
do_read() {
    request_parser_.emplace();
    request_parser_->body_limit(10000);
    stream_.expires_after(std::chrono::seconds(30));
    boost::beast::http::async_read(
        stream_,
        buffer_,
        *request_parser_,
        boost::beast::bind_front_handler(
            &http_session::on_read,
            shared_from_this()));
}

void
http_session::
on_read(boost::beast::error_code ec, std::size_t) {
    if (ec == boost::beast::http::error::end_of_stream)
        return do_close();
    if (ec)
        fail(ec, "read");
    request_handler handler(request_parser_->release());
    do_write(handler.handle_request());
}

void
http_session::
on_write(bool keep_alive, boost::beast::error_code const& ec, std::size_t) {
    if (ec) return fail(ec, "write");
    if (ec == boost::beast::http::error::end_of_stream || !keep_alive)
        return do_close();
    do_read();
}

void
http_session::
fail(boost::beast::error_code ec, std::string_view what) {
    if (ec == boost::asio::error::operation_aborted)
        return;
    std::cerr << std::format("FAIL: {}\n", what);
}

void
http_session::
do_write(boost::beast::http::response<boost::beast::http::string_body>&& response) {
    auto response_ptr { std::make_shared<boost::beast::http::response<boost::beast::http::string_body>>(response) };
    boost::beast::http::async_write(
        stream_,
        *response_ptr,
        [self = shared_from_this(), response_ptr](boost::beast::error_code const ec, std::size_t bytes) {
            self->on_write(response_ptr->keep_alive(), ec, bytes);
        }
    );
}

void
http_session::
do_close() {
    boost::beast::error_code ec;
    ec = stream_.socket().shutdown(boost::asio::ip::tcp::socket::shutdown_send, ec);
    if (ec)
        fail(ec, "close");
}

};
