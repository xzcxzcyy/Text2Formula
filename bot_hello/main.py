import sys
from PyQt5.QtWidgets import QApplication
from PyQt5.QtSvg import QSvgWidget

from io import BytesIO
import matplotlib.pyplot as plt


# matplotlib: force computer modern font set
plt.rc('mathtext', fontset='cm')


def tex2svg(formula, fontsize=52, dpi=2000):
    """Render TeX formula to SVG.
    Args:
        formula (str): TeX formula.
        fontsize (int, optional): Font size.
        dpi (int, optional): DPI.
    Returns:
        str: SVG render.
    """

    fig = plt.figure(figsize=(0.01, 0.01))
    fig.text(0, 0, r'${}$'.format(formula), fontsize=fontsize)

    output = BytesIO()
    fig.savefig(output, dpi=dpi, transparent=True, format='svg',
                bbox_inches='tight', pad_inches=0.0, frameon=False)
    plt.close(fig)
    output.seek(0)
    return output.read()


def main():
    FORMULA = sys.argv[2]
    path = sys.argv[1]
    f = open(path, 'wb')
    f.write(tex2svg(FORMULA))
    f.close()

if __name__ == '__main__':
    main()
