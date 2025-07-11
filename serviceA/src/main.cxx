#include <cstdlib>
#include <iostream>
#include "../include/listener.hxx"
#include <memory>
#include <thread>

int main(int argc, char* argv[]) {
    if (argc != 4) {
        std::cerr << "Usage: <address> <port> <threads>\n";
        return EXIT_FAILURE;
    }
    auto const address = boost::asio::ip::make_address_v4(argv[1]);
    auto const port = static_cast<unsigned short>(std::atoi(argv[2]));
    auto threads = std::atoi(argv[3]);

    boost::asio::io_context ioc(threads);

    std::make_shared<service_a::listener>(
        ioc,
        boost::asio::ip::tcp::endpoint(address, port))->run();

    boost::asio::signal_set signals(ioc, SIGINT, SIGTERM);
    signals.async_wait(
        [&ioc](boost::system::error_code const& error, int) {
            if (!error)
                std::cout << "The program was completed without errors!\n";
            else
                std::cerr << error.what() << '\n';
            ioc.stop();
        }
    );

    std::vector<std::thread> vec;
    vec.reserve(threads - 1);
    for (; threads > 0; --threads)
        vec.emplace_back([&ioc](){ioc.run();});
    for (auto& t : vec)
        t.join();
    ioc.run();

    return EXIT_SUCCESS;
}
