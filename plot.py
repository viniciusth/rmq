import pandas as pd
import matplotlib.pyplot as plt
import sys
import os
import matplotlib.ticker as ticker
import seaborn as sns
import numpy as np

def format_time_ns(y, _):
    if y == 0:
        return '0'
    if y >= 1e9:
        return f'{y/1e9:g}s'
    if y >= 1e6:
        return f'{y/1e6:g}ms'
    if y >= 1e3:
        return f'{y/1e3:g}μs'
    return f'{y:g}ns'

def format_mem_bytes(y, _):
    if y == 0:
        return '0'
    if y >= 1024**3:
        return f'{y/1024**3:g}GB'
    if y >= 1024**2:
        return f'{y/1024**2:g}MB'
    if y >= 1024:
        return f'{y/1024:g}KB'
    return f'{y:g}B'

def main():
    if len(sys.argv) < 2:
        print("Usage: python plot.py data.csv")
        sys.exit(1)
    df = pd.read_csv(sys.argv[1])
    os.makedirs("plots", exist_ok=True)
    sns.set_theme(style="darkgrid")

    unique_algos = sorted(df['algo'].unique())
    palette = dict(zip(unique_algos, plt.colormaps.get_cmap('flare')(np.linspace(0.1, 0.9, len(unique_algos)))))

    # Vary N, fixed Q
    fixed_q = 1000000
    df_n = df[(df['Q'] == fixed_q)]

    power_of_10_formatter = ticker.FuncFormatter(lambda x, _: f'$10^{{{int(np.log10(x))}}}$' if x > 0 else '0')

    # Construct time vs N
    plt.figure(figsize=(10, 6))
    ax = sns.lineplot(data=df_n, x='N', y='construct_time_ns', hue='algo', palette=palette, errorbar='sd')
    ax.set_xscale('log')
    ax.set_yscale('log')
    ax.set_xlabel('Number of Elements (N)')
    ax.set_ylabel('Construction Time')
    ax.set_title(f'Construction Time vs N (Q={fixed_q:,})')
    ax.xaxis.set_major_formatter(power_of_10_formatter)
    ax.yaxis.set_major_formatter(ticker.FuncFormatter(format_time_ns))
    plt.savefig('plots/construct_time_vs_n.png')
    plt.close()

    # Query time vs N
    plt.figure(figsize=(10, 6))
    ax = sns.lineplot(data=df_n, x='N', y='query_time_ns', hue='algo', palette=palette, errorbar='sd')
    ax.set_xscale('log')
    ax.set_yscale('log')
    ax.set_xlabel('Number of Elements (N)')
    ax.set_ylabel(f'Total Query Time for {fixed_q:,} queries')
    ax.set_title(f'Total Query Time vs N (Q={fixed_q:,})')
    ax.xaxis.set_major_formatter(power_of_10_formatter)
    ax.yaxis.set_major_formatter(ticker.FuncFormatter(format_time_ns))
    plt.savefig('plots/query_time_vs_n.png')
    plt.close()

    # Memory alloc vs N
    plt.figure(figsize=(10, 6))
    ax = sns.lineplot(data=df_n, x='N', y='construct_alloc_bytes', hue='algo', palette=palette, errorbar='sd')
    ax.set_xscale('log')
    ax.set_yscale('log', base=2)
    ax.set_xlabel('Number of Elements (N)')
    ax.set_ylabel('Memory Used')
    ax.set_title(f'Memory Used vs N (Q={fixed_q:,})')
    ax.xaxis.set_major_formatter(power_of_10_formatter)
    ax.yaxis.set_major_formatter(ticker.FuncFormatter(format_mem_bytes))
    plt.savefig('plots/memory_vs_n.png')
    plt.close()

    # Peak memory vs N
    plt.figure(figsize=(10, 6))
    ax = sns.lineplot(data=df_n, x='N', y='construct_peak_bytes', hue='algo', palette=palette, errorbar='sd')
    ax.set_xscale('log')
    ax.set_yscale('log', base=2)
    ax.set_xlabel('Number of Elements (N)')
    ax.set_ylabel('Peak Memory')
    ax.set_title(f'Peak Memory vs N (Q={fixed_q:,})')
    ax.xaxis.set_major_formatter(power_of_10_formatter)
    ax.yaxis.set_major_formatter(ticker.FuncFormatter(format_mem_bytes))
    plt.savefig('plots/peak_memory_vs_n.png')
    plt.close()

    # Vary Q, fixed N
    fixed_n = 1000000
    df_q = df[(df['N'] == fixed_n)]

    # Query time vs Q
    plt.figure(figsize=(10, 6))
    ax = sns.lineplot(data=df_q, x='Q', y='query_time_ns', hue='algo', palette=palette, errorbar='sd')
    ax.set_xscale('log')
    ax.set_yscale('log')
    ax.set_xlabel('Number of Queries (Q)')
    ax.set_ylabel('Total Query Time for Q queries')
    ax.set_title(f'Total Query Time vs Q (N={fixed_n:,})')
    ax.xaxis.set_major_formatter(power_of_10_formatter)
    ax.yaxis.set_major_formatter(ticker.FuncFormatter(format_time_ns))
    plt.savefig('plots/query_time_vs_q.png')
    plt.close()

if __name__ == "__main__":
    main() 