import { configure } from 'enzyme';
import chai from 'chai';
import sinonChai from 'sinon-chai';
import Adapter from 'enzyme-adapter-react-16';

configure({ adapter: new Adapter() });

before(() => {
  chai.should();
  chai.use(sinonChai);
});
