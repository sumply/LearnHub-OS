#include <boost/beast.hpp>
#include <boost/asio.hpp>
#include <boost/beast/core/error.hpp>
#include <boost/beast/core/tcp_stream.hpp>
#include <boost/beast/http/string_body_fwd.hpp>
#include <memory>

namespace service_a {

class http_session : public std::enable_shared_from_this<http_session> {
public:
    http_session(boost::asio::ip::tcp::socket socket);
    void run();
private:
    void fail(boost::beast::error_code ec, std::string_view what);
    void do_read();
    void on_read(boost::beast::error_code ec, std::size_t bytes_transferred);
    void do_write(boost::beast::http::response<boost::beast::http::string_body>&& response);
    void on_write(boost::beast::error_code const& ec, std::size_t bytes_transferred);
    void do_close();

    boost::beast::tcp_stream stream_;
    std::optional<boost::beast::http::response_parser<boost::beast::http::string_body>> request_parser_;
    boost::beast::flat_buffer buffer_;
};

}
