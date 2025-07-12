#pragma once

#ifdef DEBUG
    #ifndef LOG_REQUEST
        #include <iostream>
    #endif
    #define DEBUG_LOG(x) std::cerr << "[DEBUG]" << x << '\n'
    #define DEBUG_FUNC() std::cerr << "[DEBUG]" << __FILE__ << ':' << __LINE__ << " in " << __PRETTY_FUNCTION__ << '\n'
#else
    #define DEBUG_LOG(x) do {} while (0)
    #define DEBUG_FUNC() do {} while (0)
#endif

#ifdef LOG_REQUEST
    #ifndef DEBUG
        #include <iostream>
    #endif
    #define LOG_REQUEST(method, target, body) std::cerr << "[REQ]" << x << '\n';
#else
    #define LOG_REQUEST(x) do {} while (0)
#endif
