# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This script renders rating distributions. The data must be in the
# data/ratingdist/ directory and was generated using the
# rating_dist_parse.go script, which is required to run in advance.
#
# Usage:
#
# $ python rating_dist_draw.py

# Configure rendering figure to fit acmart fonts
import glob
import joypy
import os
import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import matplotlib
matplotlib.rcParams['text.usetex'] = True
matplotlib.rcParams['text.latex.preamble'] = r'''\usepackage{libertine}'''
matplotlib.rcParams['ps.usedistiller'] = 'xpdf'
sns.set_theme(style='ticks')


def draw(df, path, color):
    _, axes = joypy.joyplot(df[:240],
                            by='iteration',
                            column='rating',
                            labels=[y if y % 3 == 0 else None for y in list(
                                df[:240].iteration.unique())],
                            x_range=np.arange(-2, 8),
                            range_style='own',
                            tails=0.9,
                            grid='y',
                            linewidth=1,
                            legend=False,
                            fade=True,
                            kind='kde',
                            figsize=(3, 5),
                            color=color)
    axes[-1].set_xticks([-2, -1, 0, 1, 2, 3, 4, 5, 6, 7])
    axes[-1].set_xticklabels(['', '', '0', '1', '2', '3', '4', '5', '', ''])
    sns.despine(trim=False, left=False, top=False, right=False,)
    plt.savefig(path, bbox_inches='tight', dpi=50)


def main():
    for f in glob.glob("data/ratingdist/*.csv"):
        df = pd.read_csv(f)
        o = f.replace(".csv", ".pdf")
        o = o.replace('data', '../assets')
        draw(df, o, sns.mpl_palette('Blues', n_colors=len(df)/4))

    # Cherry-pick models to appear in the paper
    ids = ['ideal2',
           'dacdd183-4dd0-11ec-86eb-a85e4557a9b6',
           '48f9c537-43f4-4a22-a7ff-4ce8d9e48d4e',
           '086d8b7a-fa86-4b80-93b7-1a40fed9bc80']

    for index, id in enumerate(ids):
        os.system(
            f"cp ../assets/ratingdist/{id}.pdf ../assets/fig5_{index}.pdf")


if __name__ == '__main_':
    main()
