# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This scripts converts fbx files in the sessions folder to corresponding
# obj files.
#
# This script must be executed inside the dataset folder.
#
# Usage: blender -b -P convert.py

import bpy
import glob
import os


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

    # For simplicity, let's merge all mesh objects.
    scene = bpy.context.scene
    obs = []
    for ob in scene.objects:
        if ob.type == 'MESH':
            obs.append(ob)
    ctx = bpy.context.copy()
    ctx['active_object'] = obs[0]
    ctx['selected_objects'] = obs
    ctx['selected_editable_objects'] = obs
    bpy.ops.object.join(ctx)

    # export as obj file.
    bpy.ops.export_scene.obj(
        filepath=output, axis_forward='-Z', axis_up='Y', use_materials=False)


def main():
    for f in glob.glob('sessions/**/*.fbx'):
        print('start: ', f)
        i = f
        o = f.replace('fbx', 'obj')

        # skip files that already have an obj file.
        if os.path.exists(o):
            print('obj file already existed, skip: ', o)
            continue

        # these files cannot be converted to obj file using blender.
        if '2c45d01d-509f-11ec-a7cf-a85e4557a9b6' in o:
            continue
        if '94818102-4ddd-11ec-86eb-a85e4557a9b6' in o:
            continue

        # use try to prevent panic during the conversion.
        try:
            init_scene()
            fbx2obj(i, o)
        except:
            print('import/conversion error! skip: ', f)


if __name__ == '__main__':
    main()
