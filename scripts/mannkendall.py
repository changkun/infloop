# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts reproduces the results in "Section 5.1 Human-AI Mutual Interventions"
# and computes the Mann-Kendall and Augmented Dickey-Fuller to check the stationarity
# of the rating distribution.
#
# Usage:
#
# $ python mannkendall.py | tee mannkendall.txt

import numpy as np
import pandas as pd
import pymannkendall as mk
from statsmodels.tsa.stattools import adfuller


def main():
    stationary_skipped = 0
    stationary_MOS = 0

    mean_increasing = 0
    mean_decreasing = 0
    mean_notrend = 0
    var_increasing = 0
    var_decreasing = 0
    var_notrend = 0
    red_increasing = 0
    red_decreasing = 0
    red_notrend = 0

    for i in range(199):
        seq = f'data/sequence/lab/{i}.csv'
        series = pd.read_csv(seq, header=0, index_col=0).squeeze("columns")
        X = series.values

        scores = X[:, :4]
        m = np.array([np.mean(scores, axis=1)]).T
        u = np.array([np.var(scores, axis=1)]).T
        o = np.array([X[:, 4]]).T
        X = np.concatenate((m, u, o), axis=1)

        if len(X[:, 0]) >= 4:
            print(X[:, 0])
            ret = adfuller(X[:, 0])
            if ret[1] < 0.05:
                stationary_MOS += 1
        if len(X[:, 0]) < 4:
            stationary_skipped += 1

        if len(X[:, 0]) > 1:
            ret = mk.original_test(X[:, 0])
            if ret.trend == 'increasing':
                mean_increasing += 1
            if ret.trend == 'decreasing':
                mean_decreasing += 1
            if ret.trend == 'no trend':
                mean_notrend += 1

        if len(X[:, 0]) > 1:
            ret = mk.original_test(X[:, 1])
            if ret.trend == 'increasing':
                var_increasing += 1
            if ret.trend == 'decreasing':
                var_decreasing += 1
            if ret.trend == 'no trend':
                var_notrend += 1

        if len(X[1:, 2]) > 1:
            ret = mk.original_test(X[1:, 2])
            if ret.trend == 'increasing':
                red_increasing += 1
            if ret.trend == 'decreasing':
                red_decreasing += 1
            if ret.trend == 'no trend':
                red_notrend += 1

    print('stationary: ', stationary_skipped, stationary_MOS, 200 -
          stationary_skipped, 200 - stationary_skipped - stationary_MOS)
    print('mean dec trend: ', mean_decreasing, mean_notrend, mean_increasing)
    print('var dec trend: ', var_decreasing, var_notrend, var_increasing)
    print('optimal dec trend: ', red_decreasing, red_notrend, red_increasing)


if __name__ == '__main__':
    main()
