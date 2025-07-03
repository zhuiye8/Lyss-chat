(venv) root@DESKTOP-F2PJVJI:~/work/Lyss/backend# uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
INFO:     Will watch for changes in these directories: ['/root/work/Lyss/backend']
INFO:     Uvicorn running on http://0.0.0.0:8000 (Press CTRL+C to quit)
INFO:     Started reloader process [93273] using WatchFiles
Process SpawnProcess-1:
Traceback (most recent call last):
  File "/root/miniconda3/lib/python3.11/multiprocessing/process.py", line 314, in _bootstrap
    self.run()
  File "/root/miniconda3/lib/python3.11/multiprocessing/process.py", line 108, in run
    self._target(*self._args, **self._kwargs)
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/_subprocess.py", line 80, in subprocess_started
    target(sockets=sockets)
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/server.py", line 67, in run
    return asyncio.run(self.serve(sockets=sockets))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/miniconda3/lib/python3.11/asyncio/runners.py", line 190, in run
    return runner.run(main)
           ^^^^^^^^^^^^^^^^
  File "/root/miniconda3/lib/python3.11/asyncio/runners.py", line 118, in run
    return self._loop.run_until_complete(task)
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "uvloop/loop.pyx", line 1518, in uvloop.loop.Loop.run_until_complete
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/server.py", line 71, in serve
    await self._serve(sockets)
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/server.py", line 78, in _serve
    config.load()
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/config.py", line 436, in load
    self.loaded_app = import_from_string(self.app)
                      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/uvicorn/importer.py", line 19, in import_from_string
    module = importlib.import_module(module_str)
             ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/miniconda3/lib/python3.11/importlib/__init__.py", line 126, in import_module
    return _bootstrap._gcd_import(name[level:], package, level)
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "<frozen importlib._bootstrap>", line 1204, in _gcd_import
  File "<frozen importlib._bootstrap>", line 1176, in _find_and_load
  File "<frozen importlib._bootstrap>", line 1147, in _find_and_load_unlocked
  File "<frozen importlib._bootstrap>", line 690, in _load_unlocked
  File "<frozen importlib._bootstrap_external>", line 940, in exec_module
  File "<frozen importlib._bootstrap>", line 241, in _call_with_frames_removed
  File "/root/work/Lyss/backend/app/main.py", line 15, in <module>
    from app.core.config import settings
  File "/root/work/Lyss/backend/app/core/__init__.py", line 3, in <module>
    from .config import settings, get_settings
  File "/root/work/Lyss/backend/app/core/config.py", line 114, in <module>
    settings = Settings()
               ^^^^^^^^^^
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/pydantic_settings/main.py", line 188, in __init__
    super().__init__(
  File "/root/work/Lyss/backend/venv/lib/python3.11/site-packages/pydantic/main.py", line 253, in __init__
    validated_self = self.__pydantic_validator__.validate_python(data, self_instance=self)
                     ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/work/Lyss/backend/app/core/config.py", line 59, in assemble_db_connection
    username=values.get("POSTGRES_USER"),
             ^^^^^^^^^^
AttributeError: 'pydantic_core._pydantic_core.ValidationInfo' object has no attribute 'get'
^CINFO:     Stopping reloader process [93273]
(venv) root@DESKTOP-F2PJVJI:~/work/Lyss/backend# 