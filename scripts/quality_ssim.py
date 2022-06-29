# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

import pyvista as pv
from skimage.metrics import structural_similarity as ssim
import matplotlib.pyplot as plt
import torch
import pytorch3d
from pytorch3d.io import load_obj
from pytorch3d.structures import Meshes
from pytorch3d.renderer import (
    look_at_view_transform,
    FoVPerspectiveCameras,
    PointLights,
    RasterizationSettings,
    MeshRenderer,
    MeshRasterizer,
    HardFlatShader,
    BlendParams,
    Textures
)


class Render():
    def __init__(self, device: torch.device) -> None:
        """initialize a torch3d rasterization renderer.
        """
        R, T = look_at_view_transform(2, 0, 0)
        self.device = device
        self.camera = FoVPerspectiveCameras(
            znear=0.01, zfar=1000, R=R, T=T, device=self.device)
        self.renderer = MeshRenderer(
            rasterizer=MeshRasterizer(
                cameras=self.camera,
                raster_settings=RasterizationSettings(
                    perspective_correct=False,
                    image_size=512,
                    blur_radius=0.001,
                    faces_per_pixel=10,
                    bin_size=0,
                ),
            ),
            shader=HardFlatShader(
                device=self.device,
                cameras=self.camera,
                lights=PointLights(device=self.device,
                                   location=[[1.0, 1.0, 1.0]]),
                blend_params=BlendParams(background_color=(0, 0, 0)),
            )
        ).to(self.device)

    def render(self, mesh: Meshes, camera: FoVPerspectiveCameras = None) -> torch.Tensor:
        return self.renderer(mesh, cameras=camera or self.camera)


def load_and_uniform(model_path: str, device: torch.device) -> Meshes:
    """
    Load a mesh and uniform the vertex into a unit cube.
    """

    # load target mesh
    verts, faces, _ = load_obj(model_path)
    verts = verts.to(device)
    faces = faces.verts_idx.to(device)

    # rescale to the unit AABB and construct the target mesh.
    T = verts.mean(0)
    verts = verts - T
    S = max(verts.abs().max(0)[0])
    verts = verts / (1.5*S)
    return Meshes(
        verts=[verts], faces=[faces],
        textures=Textures(verts_rgb=torch.tensor(
            [0, 0.5, 1]).repeat(verts.shape[0], 1)[None].to(device)))


def calculate_reduction_ratio(src: str, dst: str) -> float:
    """calculate the reduction ratio of a mesh.

    Returns:
        float: the reduction ratio.
    """
    src_mesh = pv.read(src)
    dst_mesh = pv.read(dst)
    return (src_mesh.n_faces - dst_mesh.n_faces) / src_mesh.n_faces


def save_fig(fname: str, img: torch.Tensor):
    """save a figure. Only used for debugging the renderer."""
    plt.imshow(img)
    plt.grid('off')
    plt.axis('off')
    plt.gcf().set_facecolor('black')
    plt.savefig(fname)


def calculate_image_diff(origin, reduced, r: Render) -> float:
    """calculate the rendered image difference between two meshes. Each
    mesh is rendered from multiple viewports.

    Returns:
        float: the average image difference.
    """
    # https://scikit-image.org/docs/stable/api/skimage.metrics.html#skimage.metrics.structural_similarity

    num_views = 5
    elev = torch.linspace(0, 360, num_views)
    azim = torch.linspace(-180, 180, num_views)
    R, T = look_at_view_transform(dist=2, elev=0, azim=azim)
    target_cameras = [FoVPerspectiveCameras(
        device=r.device, R=R[None, i, ...], T=T[None, i, ...]) for i in range(num_views)]

    src_mesh = load_and_uniform(origin, r.device)
    dst_mesh = load_and_uniform(reduced, r.device)

    sum_ssim = 0.0
    for (_, cam) in enumerate(target_cameras):
        src = r.render(src_mesh, cam).cpu().detach().numpy()[0, :, :, 0:3]
        dst = r.render(dst_mesh, cam).cpu().detach().numpy()[0, :, :, 0:3]
        sum_ssim += ssim(src, dst,
                         data_range=src.max() - src.min(), channel_axis=2)

    avg_ssim = sum_ssim / num_views
    return avg_ssim


def process(device: str, model: str = 'rose'):
    """processes a given model. The model is assumed to be stored in the
    ../dataset/reduces folder. The result are written as a csv file to the
    ./data/curve folder.
    """

    with open(f'data/curve/{model}.csv', 'a') as f:
        f.write('reduce,ssim\n')
        for i in range(1, 100):
            r = Render(torch.device(device))

            src = f'../dataset/reduces/{model}/{model}_0.obj'
            dst = f'../dataset/reduces/{model}/{model}_{i}.obj'
            ratio = calculate_reduction_ratio(src, dst)
            ssim_val = calculate_image_diff(src, dst, r)
            print(src, dst, ratio, ssim_val)
            f.write(f'{ratio},{ssim_val}\n')


def main():
    device = 'cuda:0' if torch.cuda.is_available() else 'cpu'
    print(f'torch: {torch.__version__}')
    print(f'torch3d: {pytorch3d.__version__}')
    print(f'device: {device}')

    models = ['monkey', 'teapot', 'rose', 'cow', 'pumpkin']
    for model in models:
        process(device, model)


if __name__ == '__main__':
    main()
