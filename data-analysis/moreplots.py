import seaborn as sns
import matplotlib.pyplot as plt

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


def make_file_path(filename):
    return 'generated/other/' + filename
