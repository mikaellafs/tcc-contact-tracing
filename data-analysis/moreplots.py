import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from datetime import timedelta

def plot_cdf(df, column, title="", filename="cdf2.png", x=""):
    X = df[column]
    # import ipdb;ipdb.set_trace()
    sns.kdeplot(data = X, cumulative = True,cut=0)

    # plt.title(title)
    plt.xlabel(x)
    plt.ylabel("Probabilidade")
    plt.grid(True)
    plt.savefig(make_file_path(filename))
    plt.show()

def scatter_for_contact(df, title, filename="scatter.png"):
    sns.scatterplot(data=df, x="distance", y="duration")

    plt.xlabel("Distância (cm)")
    plt.ylabel("Duração (min)")

    plt.title(title)
    plt.savefig(make_file_path(filename))
    plt.show()

def contacts_by_time(df, filename="temporal.png"):
    df["date"] = df.firstcontacttimestamp.dt.date
    grouped_by_day = df.groupby("date").count()
    grouped_by_day = grouped_by_day.reset_index()

    grouped_by_day = gen_days_zero(grouped_by_day).sort_values(['date'])

    plt.xlabel("Data")
    plt.ylabel("Quantidade de registros de contato")

    print(grouped_by_day)
    plt.plot(grouped_by_day["date"], grouped_by_day["id"])
    plt.scatter(grouped_by_day["date"], grouped_by_day["id"])
    plt.savefig(make_file_path(filename))
    plt.grid(True)
    plt.show()

def make_file_path(filename):
    return 'generated/other/' + filename

def gen_days_zero(df):
    min_day = df['date'].min()
    max_day = df['date'].max()
    total_days = max_day - min_day
    print(min_day, max_day)

    for shift in range(0, total_days.days):
        current_day = min_day + timedelta(days=shift)

        if df.isin([current_day]).any().any():
            continue

        print(df.shape[1], df.columns)
        day_zero = pd.DataFrame([np.zeros(df.shape[1])], columns=df.columns)
        day_zero["date"] = current_day
        print(day_zero)

        df = pd.concat([df, day_zero])
    
    return df