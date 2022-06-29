# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts reproduces the results in "Section 5.1 Human-AI Mutual Interventions"
# and computes the Mann-Whiteney-U test to measure the difference between two rating
# populations. In addition, it renders the Figure 4b in the paper.
#
# Usage:
#
# $ python mannwhiteney.py | tee mannwhiteney.txt

from typing import Tuple
import matplotlib
import numpy as np
import pingouin as pg
import statannotations.Annotator as sa
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
sns.set_theme(style='ticks')

# Configure rendering figure to fit acmart font ``libertine''.
matplotlib.rcParams['text.usetex'] = True
matplotlib.rcParams['text.latex.preamble'] = r'''\usepackage{libertine}'''
matplotlib.rcParams['ps.usedistiller'] = 'xpdf'


class Distribution:
    def __init__(self) -> None:
        self.x = 'rating'
        self.y = 'reduction_ratio'
        self.labels = ['skip', 'terrible', 'poor', 'fair', 'good', 'excellent']

    def get_rating(self, df: pd.DataFrame, rating: str):
        return df[df['rating'] == rating]['reduction_ratio'].to_numpy()

    def get_mwu(self, df: pd.DataFrame, pair: Tuple[str, str]):
        dfa = self.get_rating(df, pair[0])
        print(np.mean(dfa))
        dfb = self.get_rating(df, pair[1])
        return pg.mwu(dfa, dfb, alternative='two-sided')

    def draw_df(self, df, ax, pairs, yticklabels=[], xlabel=None, ylabel=None):
        # 1. Major box plot.
        sns.boxplot(ax=ax, x=self.x, y=self.y, data=df,
                    palette=sns.mpl_palette('Blues'), order=self.labels)

        # 2. Annotate significants.
        anno = sa.Annotator(ax, pairs, data=df, x=self.x,
                            y=self.y, order=self.labels)
        anno.configure(test='Mann-Whitney', text_format='star', loc='inside')
        anno.apply_and_annotate()

        # 3. Visualize more data samples. Currently it is commented out because too distracting.
        # sns.stripplot(ax=ax, x='utility', y='total_reduce', data=df, size=4, color='.3', linewidth=0, order=order)

        ax.set_xticklabels(self.labels, fontsize=10, rotation=30)
        ax.set_yticklabels(yticklabels)
        ax.set(xlabel=xlabel, ylabel=ylabel)


def main():
    d = Distribution()
    df1 = pd.read_csv('data/reduce/field.csv')
    df2 = pd.read_csv('data/reduce/lab.csv')

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(6, 5))
    d.draw_df(df1, ax1, [
        # ('skip', 'terrible'),
        # ('skip', 'poor'),
        # ('terrible', 'poor'),
        # ('good', 'excellent'),
        ('fair', 'good'),
        ('fair', 'excellent'),
        ('terrible', 'excellent'),
    ], yticklabels=['', '0', '20', '40', '60', '80', '100', ''], ylabel='overall reduction ratio (\%)')
    d.draw_df(df2, ax2, [
        # ('skip', 'terrible'),
        # ('skip', 'poor'),
        # ('terrible', 'poor'),
        ('good', 'excellent'),
        ('fair', 'excellent'),
        ('terrible', 'excellent'),
    ])

    plt.savefig('../assets/fig4b.pdf', bbox_inches='tight')

    pairs = [
        ('good', 'excellent'),
        ('fair', 'good'),
        ('fair', 'excellent'),
        ('terrible', 'poor'),
        ('terrible', 'excellent'),
    ]

    for (k, v) in pairs:
        print(k, v)
        print(d.get_mwu(df1, (k, v)))
        print('median: ', np.median(d.get_rating(df1, k)),
              np.median(d.get_rating(df1, v)))
        print('mean: ', np.mean(d.get_rating(df1, k)),
              np.mean(d.get_rating(df1, v)))
        print('----')

    pairs = [
        ('good', 'excellent'),
        ('fair', 'excellent'),
        ('terrible', 'poor'),
        ('terrible', 'excellent'),
    ]

    for (k, v) in pairs:
        print(k, v)
        print(d.get_mwu(df2, (k, v)))
        print('median: ', np.median(d.get_rating(df2, k)),
              np.median(d.get_rating(df2, v)))
        print('mean: ', np.mean(d.get_rating(df2, k)),
              np.mean(d.get_rating(df2, v)))
        print('----')


if __name__ == '__main__':
    main()
