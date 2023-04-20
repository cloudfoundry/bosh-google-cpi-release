require 'rspec'
require 'json'
require 'bosh/template/test'

describe 'google_cpi job' do
  let(:release) { Bosh::Template::Test::ReleaseDir.new(File.join(File.dirname(__FILE__), '../..')) }
  let(:job) { release.job('google_cpi') }

  describe 'cpi.json' do
    let(:template) { job.template('config/cpi.json') }

    let(:config) { JSON.parse(template.render(manifest_properties)) }

    let(:manifest_properties) do
      {
        'google' => {
          'project' => 'some_google_project'
        }
      }
    end

    let(:rendered_google_properties) { config['cloud']['properties']['google'] }

    it 'renders the CPI config properly' do
      expect(rendered_google_properties['project']).to eq('some_google_project')
    end
  end
end
