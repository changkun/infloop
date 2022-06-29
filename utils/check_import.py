# Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
# Created by Changkun Ou <https://changkun.de>.
#
# Use of this source code is governed by a GNU GPLv3 license that
# can be found in the LICENSE file.

# This script checks whether a given root model is broken or not. It helps
# us identify broken models that cannot be imported by blender.
#
# The approach is to use Blender to import a root .FBX file, i.e.
# models that its model id equals to the session id. if the import
# process show any error, then we say this is a broken file.
#
# To use this script:
#
# $ blender -b -P check_import.py -- <folder>

import bpy
import os
import sys


def init_scene():
    bpy.ops.wm.read_homefile()
    bpy.ops.object.select_all(action='SELECT')
    bpy.ops.object.delete()


def pass_import(folder: str, session_id: str) -> bool:
    filepath = f'./{folder}/{session_id}/{session_id}.fbx'
    print('importing: ', filepath)

    try:
        bpy.ops.import_scene.fbx(filepath=filepath)
        return True
    except:
        print("import error!")
        return False


def args() -> str:
    """returns desired folder.
    """
    argv = sys.argv
    argv = argv[argv.index("--") + 1:]  # get all args after "--"
    return argv[0]


def main():
    folder = args()

    ids = [x[0][len(f'./{folder}/'):] for x in os.walk(f'./{folder}')]
    for sid in ids[1:]:  # skip the containing folder
        init_scene()
        if not pass_import(folder, sid):
            # write the broken model id to a file
            with open(f'broken.txt', 'a') as f:
                f.write(f'{sid}\n')
            continue


if __name__ == '__main__':
    main()
