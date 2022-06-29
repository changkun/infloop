# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts converts an fbx file to obj file
# Usage: blender --background --python fbx2obj.py -- <input_fbx_file> <output_obj_file>
import bpy
import sys


def args():
    """returns the input fbx file path and output obj file path.
    """
    argv = sys.argv
    argv = argv[argv.index("--") + 1:]  # get all args after "--"
    input = argv[0]
    output = argv[1]
    return input, output


def init_scene():
    """init initializes an empty scene for subsequent operations.
    """
    bpy.ops.wm.read_homefile()
    bpy.ops.object.select_all(action='SELECT')
    bpy.ops.object.delete()


def fbx2obj(input: str, output: str):
    """fbx2obj merges all objects, and exports to an obj file.
    """
    bpy.ops.import_scene.fbx(filepath=input, axis_forward='-Z', axis_up='Y')

    scene = bpy.context.scene
    obs = []
    for ob in scene.objects:
        # whatever objects you want to join...
        if ob.type == 'MESH':
            obs.append(ob)
    ctx = bpy.context.copy()
    ctx['active_object'] = obs[0]
    ctx['selected_objects'] = obs
    ctx['selected_editable_objects'] = obs
    bpy.ops.object.join(ctx)

    bpy.ops.export_scene.obj(
        filepath=output, axis_forward='-Z', axis_up='Y', use_materials=False)


if __name__ == "__main__":
    i, o = args()

    try:
        init_scene()
        fbx2obj(i, o)
    except Exception as e:
        print('conversion failed: ', e)
