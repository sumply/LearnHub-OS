#pragma once

#include <curl/curl.h>
#include <string>

namespace service_a {

class http_client {
public:
    enum class method { GET, POST, PUT, DELETE };

    http_client(method m, std::string_view url);
    ~http_client();
    http_client(http_client const&) = delete;
    http_client& operator=(http_client const&) = delete;

    http_client& field(std::string_view field);
    http_client& timeout(long ms);
    http_client& body(std::string_view data);
    void show();
    std::string response();
    long send();

    const char* error_message();
private:
    CURL* curl;
    curl_slist* headers;
    std::string response_data;
    std::string url;
    CURLcode error_code;

    static size_t write_callback(char* ptr, size_t size, size_t nmemb, std::string* data);

    void set_method(method m);
};

}
