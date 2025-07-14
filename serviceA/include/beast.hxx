#include <boost/beast.hpp>
#include <boost/asio.hpp>

namespace beast = boost::beast;
namespace http = beast::http;
namespace asio = boost::asio;
namespace ip = asio::ip;

using http_response = http::response<http::string_body>;
using http_request = http::request<http::string_body>;
