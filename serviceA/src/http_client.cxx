#include "../include/http_client.hxx"
#include "../include/logging.hxx"

namespace service_a {

http_client::
http_client(method m, std::string_view url)
    : curl(curl_easy_init())
    , headers(nullptr) {
        DEBUG_LOG(url);
        curl_easy_setopt(curl, CURLOPT_URL, url.data());
        set_method(m);
        curl_easy_setopt(curl, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V6);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_data);
}

http_client::
~http_client() {
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
}

http_client&
http_client::
field(std::string_view field) {
    headers = curl_slist_append(headers, field.data());
    return *this;
}

http_client&
http_client::
timeout(long ms) {
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, ms);
    return *this;
}

http_client&
http_client::
body(std::string_view data) {
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data.data());
    return *this;
}

std::string
http_client::
response() {
    return response_data;
}

long
http_client::
send() {
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    error_code = curl_easy_perform(curl);
    if (error_code != CURLE_OK)
        return -1;
    long http_code = 0;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);
    return http_code;
}

const char*
http_client::
error_message(){
    return curl_easy_strerror(error_code);
}

size_t
http_client::
write_callback(char* ptr, size_t size, size_t nmemb, std::string* data) {
    if (data) {
        data->append(ptr, size * nmemb);
        return size * nmemb;
    }
    return 0;
}

void
http_client::
set_method(method m) {
    switch (m) {
        case method::POST:
            curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
            break;
        case method::GET:
            curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
            break;
        case method::PUT:
            curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
            break;
        case method::DELETE:
            curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
            break;
    }
}

void
http_client::
show() {
    curl_easy_setopt(curl, CURLOPT_VERBOSE, 1);
    curl_easy_setopt(curl, CURLOPT_STDERR, stdout);
}

}
