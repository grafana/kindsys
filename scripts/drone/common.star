load('scripts/drone/vault.star', 'pull_secret')

ci_image = 'grafana/grafana-plugin-ci:1.6.1-alpine'

def lint():
    return {
        'name': 'go lint',
        'image': ci_image,
        'commands': [
            'make lint',
        ],
    }

def test():
    return {
        'name': 'go test',
        'image': ci_image,
        'commands': [
            'make test',
        ]
    }

def pipeline(
    name,
    trigger,
    steps,
    depends_on=[]
):

    return {
        'kind': 'pipeline',
        'type': 'docker',
        'name': name,
        'trigger': trigger,
        'steps': steps,
        'clone': {
            'retries': 3,
        },
        'volumes': [
            {
                'name': 'docker',
                'host': {
                    'path': '/var/run/docker.sock'
                }
            }
        ],
        'depends_on': depends_on,
        'image_pull_secrets': [pull_secret]
    }
