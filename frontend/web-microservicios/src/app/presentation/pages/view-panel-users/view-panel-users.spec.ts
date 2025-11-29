import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewPanelUsers } from './view-panel-users';

describe('ViewPanelUsers', () => {
  let component: ViewPanelUsers;
  let fixture: ComponentFixture<ViewPanelUsers>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewPanelUsers]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewPanelUsers);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
