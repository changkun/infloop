# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts is used to generate figure 4a in the paper.

# Configure rendering figure to fit acmart fonts
import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import matplotlib
matplotlib.rcParams['text.usetex'] = True
matplotlib.rcParams['text.latex.preamble'] = r'''\usepackage{libertine}'''
matplotlib.rcParams['ps.usedistiller'] = 'xpdf'
sns.set_theme(style='ticks')


def main():
    models = ['monkey', 'teapot', 'rose', 'cow', 'pumpkin']
    df = pd.DataFrame(columns=['reduce', 'ssim', 'model'])
    for model in models:
        df_model = pd.read_csv(f'data/curve/{model}.csv')
        df_model['model'] = model
        min = df_model['ssim'].min()
        max = df_model['ssim'].max()
        df_model['ssim'] = (df_model['ssim'] - min) / (max - min)
        df = df.append(df_model)
    df = df.reset_index()

    _, ax = plt.subplots(1, 1, figsize=(6, 5))

    ax = sns.regplot(ax=ax, x=df['reduce'], y=df['ssim'], scatter=False,
                     order=2, truncate=False, color='red', marker="None", label='regression')

    g = sns.lineplot(ax=ax, data=df, x='reduce', y='ssim', hue='model', hue_order=models,
                     style='model', ci=68, palette=sns.mpl_palette('Blues', 10)[5:])

    sns.despine(trim=False, left=False, top=False, right=False)
    g.set_xlabel('overall reduction ratio (\%)')
    g.set_ylabel('perceived visual quality (SSIM)')
    g.set_xticklabels(['0', '20', '40', '60', '80', '100'])
    g.legend_.set_title(None)

    plt.savefig('../assets/fig4a.pdf', bbox_inches='tight', dpi=50)


if __name__ == '__main__':
    main()
