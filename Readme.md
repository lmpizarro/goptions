WIP LICENSE GNU

# Options

## Calculators

[option education](https://www.optionseducation.org/toolsoptionquotes/optionscalculator)
https://golangdocs.com/golang-finance-go-package-stock-quote-options-chart


https://github.com/spzala19/Multiprocessing-with-golang

https://www.simtrade.fr/blog_simtrade/option-greeks-vega/

https://golangdocs.com/golang-finance-go-package-stock-quote-options-chart


## Implied Volatility

[Newton Raphson](https://anthonylaiuppa.com/posts/guant/)


[in julia](https://quantsrus.github.io/post/implied_volatility_algorithms_go_julia/)

[performance](https://github.com/jherekhealy/AQFED.jl)

## Formulas

[macroption](https://www.macroption.com/black-scholes-formula/)

[in python](https://medium.com/swlh/calculating-option-premiums-using-the-black-scholes-model-in-python-e9ed227afbee)

## code
[in go](https://medium.com/swlh/calculating-option-premiums-using-the-black-scholes-model-in-python-e9ed227afbee)


/*
    https://blog.quantinsti.com/implied-volatility/
	https://github.com/yassinemaaroufi/MibianLib
*/

Predicted = -392.7412 + K*2.5333 + T*0.2483 + TT*-0.0015 + KK*-0.0038 for calls
Predicted = 1007.6547 + K*-5.3601 + T*0.2808 + TT*-0.0022 + KK*0.0071 for puts

Predicted = -1801.6779 + K*9.5666 + T*0.2445 + TT*-0.0014 + KK*-0.0126 for calls
Predicted = 2488.9826 + K*-12.6720 + T*0.2784 + TT*-0.0022 + KK*0.0161 for puts

Predicted = -1801.6779 + K*9.5666 + T*0.2445 + TT*-0.0014 + KK*-0.0126 for calls
Predicted = 2488.9826 + K*-12.6720 + T*0.2784 + TT*-0.0022 + KK*0.0161 for puts

for calls Predicted = -1211.3803 + K*6.6392 + T*0.2450 + TT*-0.0014 + KK*-0.0090
for puts Predicted =1143.9836 + K*-6.0030 + T*0.2797 + TT*-0.0022 + KK*0.0079

for calls Predicted = 29415.5043 + K*-144.6617 + T*0.2199 + TT*-0.0009 + KK*0.1779
R2: 0.9146565631806203 Var Pred 8.719687264148183 Var obs 9.5332910899654
for puts Predicted = -8278.4412 + K*40.5141 + T*0.2290 + TT*-0.0012 + KK*-0.0495
R2: 0.9190303081313422 Var Pred 6.699797170520277 Var obs 7.29007205882353

volatility


# 160 for calls Predicted = 8.3918 + K*-0.0361 + T*-0.0182 + TT*-0.0000 + KK*0.0000 + KT*0.0000
R2: 0.9680544017731629 Var Pred 0.0003704772798844008 Var obs 0.0003827029547159809
# 160 for puts Predicted = -11.8856 + K*0.0645 + T*-0.0244 + TT*-0.0000 + KK*-0.0001 + KT*0.0001
R2: 0.9097037009905432 Var Pred 0.0006050517309466163 Var obs 0.0006651085735803839

# 159 for calls Predicted = 8.5276 + K*-0.0372 + T*-0.0152 + TT*-0.0000 + KK*0.0000 + KT*0.0000
R2: 9.36e-01 Var Pred 3.37e-04 Var obs 3.60e-04
# 159 for puts Predicted = -4.6745 + K*0.0281 + T*-0.0177 + TT*-0.0000 + KK*-0.0000 + KT*0.0001
R2: 8.75e-01 Var Pred 4.53e-04 Var obs 5.18e-04

# 160 for calls Predicted = 4.4469 + K*-0.0179 + T*-0.0094 + TT*-0.0000 + KK*0.0000 + KT*0.0000
R2: 9.42e-01 Var Pred 2.43e-04 Var obs 2.58e-04
# 160 for puts Predicted = -0.8526 + K*0.0078 + T*-0.0077 + TT*-0.0000 + KK*-0.0000 + KT*0.0000
R2: 9.40e-01 Var Pred 2.43e-04 Var obs 2.59e-04
