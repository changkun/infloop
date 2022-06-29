Some scripts are depending on PyTorch and PyTorch3D which are recommended to use
with a GPU acceleration. Installing PyTorch3D may be cumbersome as of June 29, 2022. This situation may change in the future.

Here we suggest a way to install PyTorch3D:

```bash
# create a conda isolated environment, namely ``infloop`.
conda create -n infloop python=3.9
conda activate infloop

# install pytorch prebuilt
pip install torch==1.11.0+cu113 torchvision==0.12.0+cu113 torchaudio==0.11.0 --extra-index-url https://download.pytorch.org/whl/cu113

# install pytorch3d dependencies
pip install fvcore iopath

# install pytorch3d prebuilt
pip install --no-index --no-cache-dir pytorch3d -f https://dl.fbaipublicfiles.com/pytorch3d/packaging/wheels/py39_cu113_pyt1110/download.html

# install all other dependencies
pip install -r requirements.txt
```

If the above approach does not work for you directly, refer to [PyTorch3D's official documentation](https://github.com/facebookresearch/pytorch3d/blob/main/INSTALL.md).