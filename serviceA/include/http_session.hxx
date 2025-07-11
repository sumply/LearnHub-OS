#include <memory>
#include "beast.hxx"

namespace service_a {

class http_session : public std::enable_shared_from_this<http_session> {
public:
    http_session(boost::asio::ip::tcp::socket socket);
    void run();
private:
    void fail(boost::beast::error_code ec, std::string_view what);
    void do_read();
    void on_read(boost::beast::error_code ec, std::size_t bytes_transferred);
    void do_write(http_response&& response);
    void on_write(bool keep_alive, boost::beast::error_code const& ec, std::size_t bytes_transferred);
    void do_close();

    boost::beast::tcp_stream stream_;
    std::optional<boost::beast::http::request_parser<boost::beast::http::string_body>> request_parser_;
    boost::beast::flat_buffer buffer_;
};

}
