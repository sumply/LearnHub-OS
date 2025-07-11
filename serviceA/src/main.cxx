#include <cstdlib>
#include <iostream>
#include "../include/listener.hxx"
#include <memory>

int main(int argc, char* argv[]) {
    if (argc != 3) {
        std::cerr << "Usage: <address> <port>\n";
        return EXIT_FAILURE;
    }
    auto const address = boost::asio::ip::make_address_v4(argv[1]);
    auto port = static_cast<unsigned short>(std::atoi(argv[2]));
    boost::asio::io_context ioc(1);

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
    ioc.run();
}
