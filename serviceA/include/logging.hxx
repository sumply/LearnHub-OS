#pragma once

#ifdef DEBUG
    #include <iostream>
    #define DEBUG_LOG(x) std::cerr << "[DEBUG]" << x << '\n'
    #define DEBUG_FUNC() std::cerr << "[DEBUG]" << __FILE__ << ':' << __LINE__ << " in " << __PRETTY_FUNCTION__ << '\n'
#else
    #define DEBUG_LOG(x) do {} while (0)
    #define DEBUG_FUNC() do {} while (0)
#endif
