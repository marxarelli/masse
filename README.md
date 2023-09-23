# Phyton

Phyton is an extensible new [BuildKit](https://docs.docker.com/build/buildkit/)
frontend that allows users to express complex container image build graphs in
YAML. It aims to:

 1. Give users simple yet powerful functional constructs to express how their
    container filesystems should be created, composed, and packaged.
 2. Provide an API for defining new build constructs.
 3. Maintain a lazy evaluation model by expressing all build instructions as
    [Low-Level Build (LLB)](https://docs.docker.com/build/buildkit/#llb).
 4. Formally separate container filesystem creation from image configuration.
 5. Give users a simple API for composing images from built filesystems and
    configuration.

## License

Phyton is licensed under the GNU General Public License 3.0 or later
(GPL-3.0+). See the LICENSE file for more details.
