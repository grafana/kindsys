load('scripts/drone/common.star', 'lint', 'test', 'pipeline')

def common_steps():
    return [
        lint(),
        test()
    ]

def pr_pipeline():
    return [
        pipeline(
            name='test-pr',
            trigger={
                'event': ['pull_request'],
            },
            steps=common_steps()
        ),
    ]

def main_pipeline():
   return [
       pipeline(
           name='test-main',
           trigger={
               'branch': ['main'],
               'event': ['push'],
           },
           steps=common_steps()
       ),
   ]
