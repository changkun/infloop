# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts reproduces the results in "Section 5.1 Human-AI Mutual Interventions"
# and computes the Kendall’s 𝜏 coefficient to measure the ordinal association
# between the reduction ratio and rating scale.
#
# Usage:
#
# $ python kendalltau.py | tee kendalltau.txt

import pandas as pd
import scipy.stats as stats


def tau(df: pd.DataFrame):
    ratio = df.reduction_ratio
    rating = df.rating.map(
        {'skip': 0, 'terrible': 1, 'poor': 2, 'fair': 3, 'good': 4, 'excellent': 5})
    tau, p_value = stats.kendalltau(ratio, rating)
    print(f'tau: {tau}, p: {p_value}')


def spearman(df: pd.DataFrame):
    ratio = df.reduction_ratio
    rating = df.rating.map(
        {'skip': 0, 'terrible': 1, 'poor': 2, 'fair': 3, 'good': 4, 'excellent': 5})
    result = stats.spearmanr(ratio, rating)
    print(result)


def main():
    field = pd.read_csv('data/reduce/field.csv')
    lab = pd.read_csv('data/reduce/lab.csv')

    tau(field)
    tau(lab)

    # While it can often be used interchangeably with Kendall’s, Kendall’s is
    # more robust and generally the preferred method of the two.
    # See: https://www.tessellationtech.io/data-science-stats-review/
    #
    # spearman(field)
    # spearman(lab)


if __name__ == '__main__':
    main()
